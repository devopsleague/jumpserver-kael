import os
import queue
import asyncio

from api import globals
from jms.base import BaseWisp
from wisp.protobuf import service_pb2
from wisp.protobuf.common_pb2 import Session, KillSession
from .session import SessionManager, SessionHandler

from utils.logger import get_logger

logger = get_logger(__name__)


class PollJMSEvent(BaseWisp):

    def clear_zombie_session(self):
        replay_dir = os.path.join(globals.PROJECT_DIR, 'data/replay')
        req = service_pb2.RemainReplayRequest(replay_dir=replay_dir)
        resp = self.stub.ScanRemainReplays(req)
        if not resp.status.ok:
            logger.error(f"Scan remain replay error: {resp.status.err}")
        else:
            print("Scan remain replay success")

    def wait_for_kill_session_message(self):
        session_handle = SessionHandler()
        resp = self.stub.DispatchTask(iter(queue.Queue(maxsize=1000).get, None))
        for task in resp:
            task_id = task.id
            session_id = task.session_id
            task_action = task.action
            target_session = None
            for session in SessionManager.get_store().values():
                if isinstance(session, Session) and session.id == session_id:
                    target_session = session
                    break
            if target_session is not None:
                if task_action == KillSession:
                    session_handle.close_session(target_session)
                req = service_pb2.FinishedTaskRequest(task_id=task_id.id)
                self.stub.FinishSession(req)

    async def start_session_killer(self):
        self.wait_for_kill_session_message()

    async def start(self):
        self.clear_zombie_session()
        await self.start_session_killer()


def setup_poll_jms_event():
    poll_jms_event = PollJMSEvent()
    asyncio.create_task(poll_jms_event.start())


if __name__ == '__main__':
    q = queue.Queue(maxsize=1000)
    poll_jms_event = PollJMSEvent()
    asyncio.create_task(poll_jms_event.start())
    print('create_task -> poll_jms_event')
    import time

    time.sleep(1)
    print('FinishedTaskRequest')
    q.put(service_pb2.FinishedTaskRequest(task_id='2'))
