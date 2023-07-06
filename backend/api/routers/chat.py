import uuid
import time
import asyncio

from typing import Optional
from fastapi import APIRouter, Depends, HTTPException
from fastapi.encoders import jsonable_encoder
from pydantic import ValidationError
from starlette.websockets import WebSocket, WebSocketDisconnect

from api.ai import ChatGPTManager
from api.message import ChatGPTMessage
from api.enums import WSStatusCode
from api.schemas import AskRequest, AskResponse, AskResponseType
from jms import SessionHandler, JMSSession

from utils import reply
from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()


# @router.get("/test")
# async def test():
#     return {"Hello": "World"}
#
#
# def operate_token1(token: Optional[str] = None):
#     if token is None:
#         raise HTTPException(status_code=400, detail="Invalid token")
#
#     # 可以在这里进行验证 token 的有效性等其他逻辑
#     token = "4acdd80c-2654-44e2-aa24-47c18f406db4"
#     token_resp = TokenHandler().sync_run('get_auth_info', token=token)
#     if not token_resp.status.ok:
#         error = token_resp.status.err
#         raise HTTPException(status_code=400, detail=error)
#
#     session_id = SessionHandler().sync_run('create', token_resp=token_resp)
#
#     d = {
#         'session_id': session_id,
#         'secret': token_resp.data.account.secret,
#         'user_name': token_resp.data.user.name,
#         'account_name': token_resp.data.account.name,
#     }
#     return d
#
#
# @router.get("/feng")
# def feng(connect_info: dict = Depends(operate_token1)):
#     return connect_info
#
#
# @router.get("/info")
# async def info():
#     from jms import CommandHandler
#     handler = CommandHandler()
#     await handler.run('test')
#
#     return {'data': 'test'}
#
#
# def create_jms_session(token: Optional[str] = None) -> JMSSession:


class BaseW:

    async def run(self):
        loop = asyncio.get_event_loop()
        loop.run_in_executor(None, self.sync_run)
        print('fengqiang')

    @staticmethod
    def sync_run():
        time.sleep(2)
        print('-------------------')


# async def create_jms_session(token: Optional[str] = None) -> JMSSession:
#     session_handler = SessionHandler()
#     jms_session = session_handler.create_new_session(token)
#     jms_session.active_session()
#     return jms_session


async def create_jms_session(token: Optional[str] = None):
    return '=---------'


def chat_func(ask_request: AskRequest, history_asks: list):
    async def inner(websocket: WebSocket):
        conversation_id = ask_request.conversation_id
        manager = ChatGPTManager()
        last_content = ''
        async for data in manager.ask(
                content=ask_request.content,
                conversation_id=conversation_id,
                history_asks=history_asks
        ):
            try:
                assert isinstance(data, ChatGPTMessage)
                message = data
                last_content = message.content
                if conversation_id is None:
                    conversation_id = uuid.uuid4()
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


@router.websocket("/chat")
# async def chat(websocket: WebSocket, jms_session: JMSSession = Depends(create_jms_session)):
async def chat(websocket: WebSocket, jms_session: str = Depends(create_jms_session)):
    print('Websocket 连接建立成功')
    await websocket.accept()

    try:
        history_asks = []
        while True:
            params = await websocket.receive_json()
            try:
                ask_request = AskRequest(**params)
            except ValidationError as e:
                logger.warning(f"Invalid ask request: {e}")
                await reply(AskResponse(type=AskResponseType.error, error_detail=str(e)))
                await websocket.close(WSStatusCode.data_error.value, "invalidAskRequest")
                return

            try:
                # await jms_session.with_audit(
                #     ask_request.content,
                #     websocket,
                #     chat_func(ask_request, history_asks)
                # )
                last_content = await chat_func(ask_request, history_asks)(websocket)
                print('------', last_content)
            except Exception as e:
                logger.error(str(e))
                await reply(
                    websocket, AskResponse(
                        type=AskResponseType.error,
                        error_detail=str(e)
                    )
                )
                await websocket.close(WSStatusCode.server_error.value, 'unknownError')
    except WebSocketDisconnect:
        logger.error('Web socket disconnect')
