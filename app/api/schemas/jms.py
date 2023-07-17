from pydantic import BaseModel
from typing import Optional

from wisp.protobuf.common_pb2 import RiskLevel


class CommandRecord(BaseModel):
    input: Optional[str] = None
    output: Optional[str] = None
    risk_level: str = RiskLevel.Normal


class JMSState(BaseModel):
    id: str
    activate_review: Optional[bool] = None
    new_dialogue: Optional[bool] = None


class Conversation(BaseModel):
    id: str
