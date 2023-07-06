import logging
import grpc

from jumpserver.chen.wisp import Service_pb2, Service_pb2_grpc
from jumpserver.chen.framework.session import Session, SessionManager
from jumpserver.chen.framework.session.impl import JMSSession
from jumpserver.chen.web.config import MockConfig

import threading

from jumpserver.chen.wisp.Service_pb2_grpc import ServiceStub
from jumpserver.chen.wisp.ServiceOuterClass import TaskResponse, FinishedTaskRequest


class RegisterJMSEvent:
    def __init__(self):
        self.service_blocking_stub = None
        self.mock_config = MockConfig()

    def set_service_blocking_stub(self, service_blocking_stub):
        self.service_blocking_stub = service_blocking_stub

    def clear_zombie_session(self):
        if self.mock_config.is_enable():
            return

        req = Service_pb2.RemainReplayRequest(replay_dir="./")
        resp = self.service_blocking_stub.scanRemainReplays(req)
        if not resp.status.ok:
            print("Scan remain replay error: {}".format(resp.status.err))
        else:
            print("Scan remain replay success")

    def start_session_killer(self):
        if self.mock_config.is_enable():
            return

        self.wait_for_kill_session_message()

    def wait_for_kill_session_message(self):
        channel = self.service_blocking_stub.channel
        stub = ServiceStub(channel)

        def dispatch_task():
            try:
                response_iterator = stub.dispatchTask(iter([]))
                for task_response in response_iterator:
                    target_session = None
                    for session in SessionManager.get_store().values():
                        if isinstance(session, JMSSession):
                            if session.get_jms_session().get_id() == task_response.task.session_id:
                                target_session = session
                                break
                    if target_session is not None:
                        if task_response.task.action == Service_pb2.TaskResponse.KillSession:
                            target_session.close()
                        req = FinishedTaskRequest(task_id=task_response.task.id)
                        request_observer.on_next(req)

            except Exception as e:
                print("waitSessionMessage error: {}".format(e))
            else:
                print("waitSessionMessage completed")

        request_observer = ServiceStub(channel).dispatchTask(iter([]))
        thread = threading.Thread(target=dispatch_task)
        thread.start()


if __name__ == '__main__':
    register_event = RegisterJMSEvent()
    with grpc.insecure_channel('wisp:50051') as channel:
        stub = Service_pb2_grpc.ServiceStub(channel)
        register_event.set_service_blocking_stub(stub)
        register_event.clear_zombie_session()
        register_event.start_session_killer()
