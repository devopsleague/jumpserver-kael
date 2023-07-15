import json
import uuid
import httpx
from typing import Optional
from datetime import datetime
from urllib.parse import urljoin

from api.conf import settings
from api.message import ChatGPTMessage, MessageType
from api.schemas import ChatGPTResponse
from utils.logger import get_logger

logger = get_logger(__name__)


def make_session(proxy: Optional[str] = None) -> httpx.AsyncClient:
    proxies = None
    if proxy:
        proxies = {'http://': proxy, 'https://': proxy}
    return httpx.AsyncClient(proxies=proxies, timeout=None)


class ChatGPTManager:

    def __init__(
            self,
            base_url: Optional[str] = None,
            api_key: Optional[str] = None,
            model: Optional[str] = None,
            proxy: Optional[str] = None,
    ):
        self.api_key = api_key
        self.ping_url = urljoin(base_url, 'models')
        self.api_url = urljoin(base_url, 'chat/completions')
        self.model = model if model else 'gpt-3.5-turbo'
        self.session = make_session(proxy)
        logger.info(f"api_url: {self.api_url} model: {model}, proxy: {proxy}")

    async def ping(self):
        try:
            response = await self.session.get(
                url=self.ping_url,
                headers={"Authorization": f"Bearer {self.api_key}"},
                timeout=httpx.Timeout(2)
            )
            return response.status_code == 200
        except Exception:
            await self.session.aclose()
            return

    async def ask(
            self, content: str, history_asks: list = None,
            extra_args: Optional[dict] = None, **_kwargs
    ):
        message_id = uuid.uuid4()
        history_asks.append(content)
        messages = history_asks[-10:]

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
        try:
            async with self.session.stream(
                    method="POST",
                    url=self.api_url,
                    json=data,
                    headers={"Authorization": f"Bearer {self.api_key}"},
                    timeout=httpx.Timeout(settings.chat_gpt.timeout)
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
                        error_message = 'ChatGPTResponse parse json error'
                        raise Exception(error_message)

                reply_message.type = MessageType.finish
                yield reply_message

        except httpx.TimeoutException:
            error_message = 'Connect Chat GPT timeout'
            raise Exception(error_message)
        except httpx.ConnectError:
            error_message = f'Connect Chat GPT error'
            raise Exception(error_message)
