import os
import asyncio

from api import globals
from jms.base import BaseWisp
from wisp.protobuf import service_pb2
from wisp.protobuf.common_pb2 import Session, KillSession
from .session import SessionManager, SessionHandler

from utils.logger import get_logger

logger = get_logger(__name__)


class PollJMSEvent(BaseWisp):
    def __init__(self):
        super().__init__()
        self.request_observer = None

    def clear_zombie_session(self):
        replay_dir = os.path.join(globals.PROJECT_DIR, 'data/replay')
        req = service_pb2.RemainReplayRequest(replay_dir=replay_dir)
        resp = self.stub.ScanRemainReplays(req)
        if not resp.status.ok:
            logger.error(f"Scan remain replay error: {resp.status.err}")
        else:
            logger.info("Scan remain replay success")

    async def on_next_task_response(self, task_response):
        target_session = None
        for session in SessionManager.get_store().values():
            if isinstance(session, Session):
                if session.id == task_response.task.session_id:
                    target_session = session
                    break

        if target_session is not None:
            if task_response.task.action == KillSession:
                SessionHandler().close_session(target_session)
            req = service_pb2.FinishedTaskRequest(task_id=task_response.task.id)
            self.request_observer.on_next(req)

    def wait_for_kill_session_message(self):
        self.request_observer = self.stub.DispatchTask(
            self.on_next_task_response
        )

    async def start_session_killer(self):
        self.wait_for_kill_session_message()

    async def start(self):
        self.clear_zombie_session()
        await self.start_session_killer()


def setup_poll_jms_event():
    poll_jms_event = PollJMSEvent()
    asyncio.run(poll_jms_event.start())
