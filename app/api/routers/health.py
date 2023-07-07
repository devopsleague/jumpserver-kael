import datetime
from fastapi import APIRouter

from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()

upTime = datetime.datetime.now()


@router.get("/health")
async def health():
    status = {}
    now = datetime.datetime.now()
    status["timestamp"] = now.utcnow()
    status["uptime"] = str(now - upTime)
    return status
