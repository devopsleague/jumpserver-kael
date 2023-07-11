import json
import uuid
from datetime import datetime
from typing import Optional

import httpx
from pydantic import ValidationError

from api.conf import settings
from api.message import ChatGPTMessage
from api.schemas import ChatGPTResponse
from utils.logger import get_logger

logger = get_logger(__name__)


def make_session(proxy: Optional[str] = None) -> httpx.AsyncClient:
    if proxy:
        proxies = {
            'http://': proxy,
            'https://': proxy,
        }
        session = httpx.AsyncClient(proxies=proxies, timeout=None)
    else:
        session = httpx.AsyncClient(timeout=None)
    return session


class ChatGPTManager:

    def __init__(
            self,
            api_key: Optional[str] = None,
            model: Optional[str] = None,
            proxy: Optional[str] = None,
    ):
        self.api_key = api_key
        self.model = model if model else 'gpt-3.5-turbo'
        self.session = make_session(proxy)

    async def ask(
            self, content: str, history_asks: list = None,
            extra_args: Optional[dict] = None, **_kwargs
    ):
        message_id = uuid.uuid4()

        history_asks.append(content)
        messages = history_asks[-10:]

        base_url = settings.chat_gpt.openai_base_url
        data = {
            "model": self.model,
            "messages": [
                {"role": 'user', "content": msg}
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
                headers={"Authorization": f"Bearer {self.api_key}"},
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
