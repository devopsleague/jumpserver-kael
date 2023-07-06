import threading

from jms.base import BaseWisp
from wisp.protobuf import service_pb2
from wisp.protobuf.common_pb2 import Session, KillSession
from .session import SessionManager, SessionHandler


class PollJMSEvent(BaseWisp):

    def clear_zombie_session(self):
        req = service_pb2.RemainReplayRequest(replay_dir="./")
        resp = self.stub.scanRemainReplays(req)
        if not resp.status.ok:
            print("Scan remain replay error: {}".format(resp.status.err))
        else:
            print("Scan remain replay success")

    def start_session_killer(self):
        self.wait_for_kill_session_message()

    def wait_for_kill_session_message(self):
        session_handler = SessionHandler()

        def dispatch_task():
            try:
                response_iterator = self.stub.DispatchTask(iter([]))
                for task_response in response_iterator:
                    target_session = None
                    for session in SessionManager.get_store().values():
                        if isinstance(session, Session):
                            if session.id == task_response.task.session_id:
                                target_session = session
                                break
                    if target_session is not None:
                        if task_response.task.action == KillSession:
                            session_handler.close_session(target_session)
                        req = service_pb2.FinishedTaskRequest(task_id=task_response.task.id)
                        request_observer.on_next(req)
            except Exception as e:
                print("waitSessionMessage error: {}".format(e))
            else:
                print("waitSessionMessage completed")

        request_observer = self.stub.DispatchTask(iter([]))
        thread = threading.Thread(target=dispatch_task)
        thread.start()
