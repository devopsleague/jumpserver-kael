import datetime
import uuid
from enum import auto
from strenum import StrEnum
from pydantic import BaseModel
from typing import Optional, Literal


class MessageType(StrEnum):
    message = auto()
    finish = auto()


class ChatGPTMessage(BaseModel):
    content: str
    id: uuid.UUID | str
    parent: Optional[uuid.UUID]
    create_time: Optional[datetime.datetime]
    type: MessageType = MessageType.message
    role: Literal['system', 'user', 'assistant'] = 'user'
