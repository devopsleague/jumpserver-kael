import asyncio
from datetime import datetime
from typing import Optional

from starlette.websockets import WebSocket

from i18n import gettext as _
from api.jms.base import BaseWisp
from api.schemas import (
    AskResponse, AskResponseType, CommandRecord, JMSState
)
from wisp.protobuf import service_pb2
from wisp.exceptions import WispError
from wisp.protobuf.common_pb2 import TokenAuthInfo, Session
from utils.logger import get_logger
from utils.ws import reply
from ..replay import ReplayHandler
from ..command import CommandHandler

logger = get_logger(__name__)


class JMSSession(BaseWisp):
    def __init__(self, session: Session, auth_info: TokenAuthInfo, websocket: WebSocket):
        super().__init__()
        self.session = session
        self.websocket = websocket
        self.history_asks = []
        self.current_ask_interrupt = False
        self.command_acls = list(auth_info.filter_rules)
        self.expire_time = auth_info.expire_info.expire_at
        self.max_idle_time_delta = auth_info.setting.max_idle_time
        self.session_handler = None
        self.command_handler = None
        self.replay_handler = None
        self.jms_state = JMSState(id=session.id)

    async def active_session(self) -> None:
        from .manager import SessionManager
        SessionManager.register_jms_session(self)
        self.replay_handler = ReplayHandler(self.session)
        self.session_handler = SessionHandler(self.websocket)
        self.command_handler = CommandHandler(
            self.websocket, self.session,
            self.command_acls, self.jms_state
        )
        asyncio.create_task(self.maximum_idle_time_detection())

    async def maximum_idle_time_detection(self):
        last_active_time = datetime.now()

        while True:
            current_time = datetime.now()
            idle_time = current_time - last_active_time

            if idle_time.total_seconds() >= self.max_idle_time_delta * 60:
                await self.close()
                break

            if self.jms_state.new_dialogue:
                last_active_time = current_time
                self.jms_state.new_dialogue = False

            await asyncio.sleep(3)

    async def close(self) -> None:
        from .manager import SessionManager
        self.current_ask_interrupt = True
        await asyncio.sleep(1)
        await self.replay_handler.upload()
        await self.session_handler.close_session(self.session)
        SessionManager.unregister_jms_session(self)
        await self.notify_to_close()

    async def notify_to_close(self):
        await reply(
            self.websocket, AskResponse(
                type=AskResponseType.finish,
                conversation_id=self.session.id,
                system_message=_('Session interrupted')
            )
        )

    async def with_audit(self, command: str, chat_func):
        command_record = CommandRecord(input=command)
        self.command_handler.command_record = command_record
        try:
            is_continue = await self.command_handler.command_acl_filter()
            asyncio.create_task(self.replay_handler.write_input(command_record.input))
            if not is_continue:
                return

            result = await chat_func(self)
            command_record.output = result
            asyncio.create_task(self.replay_handler.write_output(command_record.output))
            return result

        except Exception as e:
            error = str(e)
            asyncio.create_task(self.replay_handler.write_output(error))
            raise e

        finally:
            asyncio.create_task(self.command_handler.record_command())


class SessionHandler(BaseWisp):

    def __init__(self, websocket: Optional[WebSocket] = None):
        super().__init__()
        self.websocket = websocket
        self.remote_address = self.get_remote_address()

    def get_remote_address(self) -> str:
        websocket = self.websocket
        remote_address = websocket.client.host
        if "x-forwarded-for" in websocket.headers:
            remote_address = websocket.headers["x-forwarded-for"]
            remote_address = remote_address.split(',')[0].strip()
        return remote_address

    def create_new_session(self, auth_info: TokenAuthInfo) -> JMSSession:
        session = self.create_session(auth_info)
        return JMSSession(session, auth_info, self.websocket)

    def create_session(self, auth_info: TokenAuthInfo) -> Session:
        req_session = Session(
            user_id=auth_info.user.id,
            user=f'{auth_info.user.name}({auth_info.user.username})',
            account_id=auth_info.account.id,
            account=f'{auth_info.account.name}({auth_info.account.username})',
            org_id=auth_info.asset.org_id,
            asset_id=auth_info.asset.id,
            asset=auth_info.asset.name,
            login_from=Session.LoginFrom.WT,
            protocol=auth_info.asset.protocols[0].name,
            date_start=int(datetime.now().timestamp()),
            remote_addr=self.remote_address,
        )
        req = service_pb2.SessionCreateRequest(data=req_session)
        resp = self.stub.CreateSession(req)
        if not resp.status.ok:
            error_message = f'Failed to create session: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)
        return resp.data

    async def close_session(self, session: Session) -> None:
        req = service_pb2.SessionFinishRequest(
            id=session.id,
            date_end=int(datetime.now().timestamp())
        )
        resp = self.stub.FinishSession(req)

        if not resp.status.ok:
            error_message = f'Failed to close session: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)
