import json

from typing import Optional
from fastapi import APIRouter, Depends
from pydantic import ValidationError
from starlette import status
from starlette.responses import Response
from starlette.websockets import WebSocket, WebSocketDisconnect

from api.ai import ChatGPTManager
from api.message import ChatGPTMessage, Conversation
from api.schemas import AskRequest, AskResponse, AskResponseType
from jms import SessionHandler, JMSSession, TokenHandler, SessionManager
from wisp.protobuf.common_pb2 import TokenAuthInfo

from utils import reply
from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()


@router.post("/interrupt_current_ask")
async def interrupt_current_ask(conversation: Conversation):
    jms_session = SessionManager.get_jms_session(conversation.id)
    if jms_session:
        assert isinstance(jms_session, JMSSession)
        jms_session.current_ask_interrupt = True
        return Response(status_code=status.HTTP_200_OK)
    else:
        return Response(status_code=status.HTTP_404_NOT_FOUND)


# async def create_auth_info(token: Optional[str] = None) -> TokenAuthInfo:
#     token_handler = TokenHandler()
#     auth_info = token_handler.get_token_auth_info(token)
#     return auth_info


async def create_auth_info(token: Optional[str] = None):
    return 'test'


@router.websocket("/chat")
# async def chat(websocket: WebSocket, auth_info: TokenAuthInfo = Depends(create_auth_info)):
async def chat(websocket: WebSocket, auth_info: str = Depends(create_auth_info)):
    session_handler = SessionHandler(websocket)
    await websocket.accept()
    print('Websocket 连接建立成功')
    current_jms_sessions = []
    try:
        async for message in websocket.iter_text():
            try:
                message = json.loads(message)
            except json.JSONDecodeError:
                await websocket.send_text("pong")
                continue

            try:
                ask_request = AskRequest(**message)
            except ValidationError as e:
                logger.warning(f"Invalid ask request: {e}")
                await reply(websocket, AskResponse(type=AskResponseType.error, system_message=str(e)))
                continue

            if ask_request.conversation_id is None:
                jms_session = session_handler.create_new_session(auth_info)
                jms_session.active_session()
                current_jms_sessions.append(current_jms_sessions)
            else:
                conversation_id = ask_request.conversation_id
                jms_session = SessionManager.get_jms_session(conversation_id)
                assert isinstance(jms_session, JMSSession)

            await jms_session.with_audit(
                ask_request.content,
                chat_func(ask_request)
            )
            # await chat_func(ask_request, history_asks)(conversation_id, websocket)

    except WebSocketDisconnect as e:
        logger.error('Web socket disconnect', e)
        for jms_session in current_jms_sessions:
            SessionManager.unregister_jms_session(jms_session)


def chat_func(ask_request: AskRequest):
    manager = ChatGPTManager()

    async def inner(jms_session: JMSSession):
        websocket = jms_session.websocket
        conversation_id = jms_session.session.id
        history_asks = jms_session.history_asks
        last_content = ''
        async for data in manager.ask(
                content=ask_request.content,
                conversation_id=conversation_id,
                history_asks=history_asks
        ):

            if jms_session.current_ask_interrupt:
                jms_session.current_ask_interrupt = False
                break

            try:
                assert isinstance(data, ChatGPTMessage)
                message = data
                last_content = message.content
            except Exception as e:
                logger.warning(f"convert message error: {e}")
                continue

            await reply(
                websocket, AskResponse(
                    type=AskResponseType.message,
                    conversation_id=conversation_id,
                    message=message
                )
            )
        await reply(
            websocket, AskResponse(
                type=AskResponseType.finish,
                conversation_id=conversation_id
            )
        )
        return last_content

    return inner
