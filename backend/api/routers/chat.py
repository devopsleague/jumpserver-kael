from fastapi import APIRouter
from fastapi.encoders import jsonable_encoder
from pydantic import ValidationError
from starlette.websockets import WebSocket, WebSocketDisconnect

from api.ai import ChatGPTManager
from api.message import ChatGPTMessage
from api.enums import WSStatusCode
from api.schemas import AskRequest, AskResponse, AskResponseType

from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()


@router.websocket("/chat")
async def chat(websocket: WebSocket):

    async def reply(response: AskResponse):
        await websocket.send_json(jsonable_encoder(response))

    await websocket.accept()
    # user = await websocket_auth(websocket)

    try:
        while True:
            params = await websocket.receive_json()
            try:
                ask_request = AskRequest(**params)
            except ValidationError as e:
                logger.warning(f"Invalid ask request: {e}")
                await reply(AskResponse(type=AskResponseType.error, error_detail=str(e)))
                await websocket.close(WSStatusCode.data_error.value, "invalidAskRequest")
                return

            # 命令复核等...
            pass

            conversation_id = None
            new_conversation = ask_request.new_conversation
            if not new_conversation:
                conversation_id = ask_request.conversation_id

            try:
                await reply(AskResponse(
                    type=AskResponseType.waiting
                ))

                manager = ChatGPTManager()

                async for data in manager.ask(
                        content=ask_request.content,
                        conversation_id=conversation_id,
                        parent_id=ask_request.parent,
                ):
                    try:
                        assert isinstance(data, ChatGPTMessage)
                        message = data
                        conversation_id = message.id
                    except Exception as e:
                        logger.warning(f"convert message error: {e}")
                        continue

                    await reply(AskResponse(
                        type=AskResponseType.message,
                        conversation_id=conversation_id,
                        message=message
                    ))
            except Exception as e:
                logger.error(str(e))
                await reply(AskResponse(
                    type=AskResponseType.error,
                    error_detail=str(e)
                ))
                await websocket.close(WSStatusCode.server_error.value, 'unknownError')
    except WebSocketDisconnect:
        logger.error('Web socket disconnect')
