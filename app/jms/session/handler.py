import asyncio
from typing import Optional
from starlette.websockets import WebSocket

from datetime import datetime
from jms.base import BaseWisp

from wisp.protobuf import service_pb2
from wisp.exceptions import WispError
from wisp.protobuf.common_pb2 import TokenAuthInfo, Session
from utils.logger import get_logger
from ..replay import ReplayHandler
from ..command import CommandHandler, CommandRecord

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

    def active_session(self) -> None:
        from .manager import SessionManager
        SessionManager.register_jms_session(self)
        self.session_handler = SessionHandler(self.websocket)
        self.command_handler = CommandHandler(self.websocket, self.session, self.command_acls)
        self.replay_handler = ReplayHandler(self.session)

    def close(self) -> None:
        from .manager import SessionManager
        self.current_ask_interrupt = True
        self.replay_handler.upload()
        self.session_handler.close_session(self.session)
        SessionManager.unregister_jms_session(self)

    async def with_audit(self, command: str, chat_func):
        command_record = CommandRecord(input=command)
        try:
            is_continue = self.command_handler.command_acl_filter(command_record)
            asyncio.create_task(self.replay_handler.write_input(command_record.input))
            if not is_continue:
                return

            result = await chat_func(self)
            command_record.output = result
            asyncio.create_task(self.replay_handler.write_input(result.output))
            return result

        except Exception as e:
            error = str(e)
            asyncio.create_task(self.replay_handler.write_input(error))
            raise e

        finally:
            asyncio.create_task(self.command_handler.record_command(command_record))


class SessionHandler(BaseWisp):

    def __init__(self, websocket: Optional[WebSocket] = None):
        super().__init__()
        self.websocket = websocket

    def create_new_session(self, auth_info: TokenAuthInfo) -> JMSSession:
        session = self.create_session(auth_info)

        try:
            # TODO 要不要看一下 secret 能不能连上
            pass
        except Exception as e:
            self.close_session(session)
            raise e

        return JMSSession(session, auth_info, self.websocket)

    def create_session(self, auth_info: TokenAuthInfo) -> Session:
        req_session = Session(
            user_id=auth_info.user.id,
            user=f"{auth_info.user.name}({auth_info.user.username})",
            account_id=auth_info.account.id,
            account=f"{auth_info.account.name}({auth_info.account.username})",
            org_id=auth_info.asset.org_id,
            asset_id=auth_info.asset.id,
            asset=auth_info.asset.name,
            login_from=Session.LoginFrom.WT,
            protocol=auth_info.asset.protocols[0].name,
            date_start=int(datetime.now().timestamp())
        )
        req = service_pb2.SessionCreateRequest(data=req_session)
        resp = self.stub.CreateSession(req)
        if not resp.status.ok:
            error_message = f'Failed to create session: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)
        return resp.data

    def close_session(self, session: Session) -> None:
        req = service_pb2.SessionFinishRequest(
            id=session.id,
            date_end=int(datetime.now().timestamp())
        )
        resp = self.stub.FinishSession(req)

        if not resp.status.ok:
            error_message = f'Failed to close session: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)
