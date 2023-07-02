from typing import Literal, Optional

from pydantic import BaseModel


class ChatGPTResponseChoice(BaseModel):
    index: Optional[int]
    delta: Optional[dict[Literal["role", "content"], str]]
    finish_reason: Optional[str]


class ChatGPTResponse(BaseModel):
    id: str
    created: int
    model: str
    choices: Optional[list[ChatGPTResponseChoice]]


# {
#     "id": "chatcmpl-7XmcU6k5YIUFLY36gYW9HPtarGplX",
#     "object": "chat.completion.chunk",
#     "created": 1688286074,
#     "model": "gpt-3.5-turbo-0613",
#     "choices": [{"index": 0, "delta": {"role":"assistant", "content": "\xe4\xbd\xa0"}, "finish_reason": None}]
# }


# b'data: {"id":"chatcmpl-7XmcU6k5YIUFLY36gYW9HPtarGplX","object":"chat.completion.chunk","created":1688286074,"model":"gpt-3.5-turbo-0613","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}'
# b''
# b'data: [DONE]'
# b''