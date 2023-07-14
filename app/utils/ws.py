import typing
from pydantic import BaseModel
from fastapi.encoders import jsonable_encoder
from starlette.websockets import WebSocket


async def reply(websocket: WebSocket, response: BaseModel):
    await websocket.send_json(jsonable_encoder(response))


async def iter_text(websocket: WebSocket) -> typing.AsyncIterator[str]:
    while True:
        yield await websocket.receive_text()
