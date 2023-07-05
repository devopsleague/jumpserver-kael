import json
import uuid
from datetime import datetime
from typing import Optional

import httpx
from pydantic import ValidationError

from api.conf import settings
from api.enums import ChatGPTModels
from api.message import ChatGPTMessage
from api.schemas import ChatGPTResponse
from utils.logger import get_logger

logger = get_logger(__name__)

api_key = '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'


def make_session() -> httpx.AsyncClient:
    if settings.chat_gpt.proxy is not None:
        proxies = {
            "http://": settings.chat_gpt.proxy,
            "https://": settings.chat_gpt.proxy,
        }
        session = httpx.AsyncClient(proxies=proxies, timeout=None)
    else:
        session = httpx.AsyncClient(timeout=None)
    return session


class ChatGPTManager:

    def __init__(self):
        self.session = make_session()

    async def ask(
            self, content: str, conversation_id: uuid.UUID = None,
            parent_id: uuid.UUID = None, model: ChatGPTModels = None,
            extra_args: Optional[dict] = None, **_kwargs
    ):
        model = ChatGPTModels('gpt_3_5')
        message_id = uuid.uuid4()
        new_message = ChatGPTMessage(
            id=message_id,
            role='user',
            create_time=datetime.now(),
            content=content,
            parent=parent_id
        )

        messages = []
        if not conversation_id:
            messages = [new_message]
        else:
            # TODO 当前不记录历史
            messages.append(new_message)

        base_url = settings.chat_gpt.openai_base_url
        data = {
            "model": model.code(),
            "messages": [
                {"role": msg.role, "content": msg.content}
                for msg in messages
            ],
            "stream": True,
            **(extra_args or {})
        }

        text_content = ''
        reply_message = None

        read_timeout = settings.chat_gpt.read_timeout
        connect_timeout = settings.chat_gpt.connect_timeout
        timeout = httpx.Timeout(read_timeout, connect=connect_timeout)

        async with self.session.stream(
                method="POST",
                url=f"{base_url}chat/completions",
                json=data,
                headers={"Authorization": f"Bearer {api_key}"},
                timeout=timeout
        ) as response:
            async for line in response.aiter_lines():
                if not line or line is None:
                    continue
                if "data: " in line:
                    line = line[6:]
                if "[DONE]" in line:
                    break

                try:
                    line = json.loads(line)
                    resp = ChatGPTResponse(**line)

                    if resp.choices[0].delta is not None:
                        text_content += resp.choices[0].delta.get('content', '')
                    if reply_message is None:
                        reply_message = ChatGPTMessage(
                            id=uuid.uuid4(),
                            role='assistant',
                            create_time=datetime.now(),
                            content=text_content,
                            parent=message_id,
                        )
                    else:
                        reply_message.content = text_content

                    yield reply_message

                except json.decoder.JSONDecodeError:
                    logger.warning(f"ChatGPTResponse parse json error")
                except ValidationError as e:
                    logger.warning(f"ChatGPTResponse validate error: {e}")
