import typing
from pydantic import BaseModel
from fastapi.encoders import jsonable_encoder
from starlette.websockets import WebSocket

from utils.logger import get_logger

logger = get_logger(__name__)


async def reply(websocket: WebSocket, response: BaseModel):
    try:
        await websocket.send_json(jsonable_encoder(response))
    except Exception as e:
        logger.error(f'websocket error: {e}')


async def iter_text(websocket: WebSocket) -> typing.AsyncIterator[str]:
    while True:
        yield await websocket.receive_text()
