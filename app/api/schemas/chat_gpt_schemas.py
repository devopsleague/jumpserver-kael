from typing import Literal, Optional

from pydantic import BaseModel


class ChatGPTResponseChoice(BaseModel):
    index: Optional[int]
    delta: Optional[dict[Literal['role', 'content'], str]]
    finish_reason: Optional[str]


class ChatGPTResponse(BaseModel):
    id: str
    created: int
    model: str
    choices: Optional[list[ChatGPTResponseChoice]]
