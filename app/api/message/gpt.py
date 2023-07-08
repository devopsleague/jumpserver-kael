import datetime
import uuid
from typing import Optional, Literal

from pydantic import BaseModel


class ChatGPTMessage(BaseModel):
    id: uuid.UUID | str
    parent: Optional[uuid.UUID]
    role: Literal['system', 'user', 'assistant'] = 'user'
    create_time: Optional[datetime.datetime]
    content: str


class Conversation(BaseModel):
    id: str
