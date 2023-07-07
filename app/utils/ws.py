from fastapi.encoders import jsonable_encoder
from starlette.websockets import WebSocket
from pydantic import BaseModel


async def reply(websocket: WebSocket, response: BaseModel):
    await websocket.send_json(jsonable_encoder(response))
