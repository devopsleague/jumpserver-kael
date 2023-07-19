import json
import asyncio

from typing import Optional
from fastapi import APIRouter, Depends, HTTPException
from pydantic import ValidationError
from starlette import status
from starlette.responses import Response
from starlette.websockets import WebSocket, WebSocketDisconnect

from api.ai import ChatGPTManager
from api.message import ChatGPTMessage, MessageType
from api.schemas import AskRequest, AskResponse, AskResponseType, Conversation, JMSState
from api.jms import SessionHandler, JMSSession, TokenHandler, SessionManager
from wisp.protobuf.common_pb2 import TokenAuthInfo
from wisp.exceptions import WispError

from utils import reply
from utils.ws import iter_text
from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()


@router.post("/interrupt_current_ask/")
async def interrupt_current_ask(conversation: Conversation):
    jms_session = SessionManager.get_jms_session(conversation.id)
    if jms_session:
        assert isinstance(jms_session, JMSSession)
        jms_session.current_ask_interrupt = True
        return Response(status_code=status.HTTP_200_OK)
    else:
        return Response('Not found conversation id', status_code=status.HTTP_404_NOT_FOUND)


@router.post("/jms_state/")
async def jms_state(state: JMSState):
    jms_session = SessionManager.get_jms_session(state.id)
    if jms_session:
        assert isinstance(jms_session, JMSSession)
        jms_session.jms_state.activate_review = state.activate_review
        return Response(status_code=status.HTTP_200_OK)
    else:
        return Response('Not found conversation id', status_code=status.HTTP_404_NOT_FOUND)


async def create_auth_info(token: Optional[str] = None) -> TokenAuthInfo:
    token_handler = TokenHandler()
    try:
        auth_info = token_handler.get_token_auth_info(token)
    except WispError as e:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail=str(e))
    return auth_info


@router.websocket("/chat/")
async def chat(websocket: WebSocket, auth_info: TokenAuthInfo = Depends(create_auth_info)):
    session_handler = SessionHandler(websocket)
    await websocket.accept()
    current_jms_sessions = []
    api_key = auth_info.account.secret
    base_url = auth_info.asset.address
    proxy = auth_info.asset.specific.http_proxy
    model = auth_info.platform.protocols[0].settings.get('api_mode')
    manager = ChatGPTManager(base_url=base_url, api_key=api_key, model=model, proxy=proxy)

    if not await manager.ping():
        await websocket.close(status.WS_1008_POLICY_VIOLATION)
        return

    try:
        async for message in iter_text(websocket):
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
            try:
                if ask_request.conversation_id is None:
                    jms_session = session_handler.create_new_session(auth_info)
                    await jms_session.active_session()
                    current_jms_sessions.append(jms_session)
                else:
                    conversation_id = ask_request.conversation_id
                    jms_session = SessionManager.get_jms_session(conversation_id)
                    if jms_session is None:
                        await reply(
                            websocket,
                            AskResponse(
                                type=AskResponseType.error,
                                system_message='Not found session id'
                            )
                        )
                        continue
                    jms_session.jms_state.new_dialogue = True

                asyncio.create_task(
                    jms_session.with_audit(
                        ask_request.content,
                        chat_func(ask_request, manager)
                    )
                )
            except WispError as e:
                logger.error(e)

    except WebSocketDisconnect:
        logger.warning('Web socket disconnect')
        for jms_session in current_jms_sessions:
            await jms_session.close()


def chat_func(ask_request: AskRequest, manager: ChatGPTManager):
    async def inner(jms_session: JMSSession):
        websocket = jms_session.websocket
        history_asks = jms_session.history_asks
        conversation_id = jms_session.session.id
        last_content = ''
        interrupt = False
        try:
            async for message in manager.ask(
                    content=ask_request.content,
                    history_asks=history_asks
            ):
                assert isinstance(message, ChatGPTMessage)
                last_content = message.content

                if jms_session.current_ask_interrupt:
                    interrupt = True
                    message.type = MessageType.finish
                    jms_session.current_ask_interrupt = False

                response_message = AskResponse(
                    type=AskResponseType.message,
                    conversation_id=conversation_id,
                    message=message
                )

                await reply(websocket, response_message)

                if interrupt:
                    break
        except Exception as e:
            logger.error(f"chat error: {e}")
            await reply(
                websocket, AskResponse(
                    type=AskResponseType.error,
                    conversation_id=conversation_id,
                    system_message=str(e)
                )
            )
        return last_content

    return inner
