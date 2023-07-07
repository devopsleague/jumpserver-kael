import uuid
from enum import auto
from typing import Optional

from pydantic import BaseModel
from strenum import StrEnum

from api.message import ChatGPTMessage
from utils.logger import get_logger

logger = get_logger(__name__)


class AskRequest(BaseModel):
    conversation_id: Optional[uuid.UUID] = None
    content: str
    parent: Optional[uuid.UUID] = None


class AskResponseType(StrEnum):
    waiting = auto()
    reject = auto()
    message = auto()
    error = auto()
    finish = auto()


class AskResponse(BaseModel):
    type: AskResponseType
    conversation_id: uuid.UUID = None
    message: Optional[ChatGPTMessage] = None
    system_message: str = None
