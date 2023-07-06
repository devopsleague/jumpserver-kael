from datetime import datetime
from jms.base import BaseWisp

from wisp.protobuf import service_pb2
from wisp.protobuf.common_pb2 import TokenAuthInfo, Session

from ..replay import ReplayHandler
from ..command import CommandHandler, CommandRecord


class SessionHandler(BaseWisp):

    def create_new_session(self, token):
        token_resp = self.get_token_auth_info(token)
        jms_session = self.create_session(token_resp)

        try:
            # TODO 要不要看一下 secret 能不能连上
            pass
        except Exception as e:
            self.close_session(jms_session)
            raise e

        return JMSSession(jms_session, token_resp)

    def get_token_auth_info(self, token: str) -> TokenAuthInfo:
        req = service_pb2.TokenRequest(token=token)
        token_resp = self.stub.GetTokenAuthInfo(req)
        if not token_resp.status.ok:
            error = token_resp.status.err
            print('获取 token 失败', error)

        return token_resp.data

    def create_session(self, token_resp: TokenAuthInfo) -> Session:
        jms_session = Session(
            user_id=token_resp.user.id,
            user=f"{token_resp.user.name}({token_resp.user.username})",
            account_id=token_resp.account.id,
            account=f"{token_resp.account.name}({token_resp.account.username})",
            org_id=token_resp.asset.org_id,
            asset_id=token_resp.asset.id,
            asset=token_resp.asset.name,
            login_from=Session.LoginFrom.WT,
            protocol=token_resp.asset.protocols[0].name,
            date_start=int(datetime.now().timestamp())
        )
        req = service_pb2.SessionCreateRequest(data=jms_session)
        resp = self.stub.CreateSession(req)
        if not resp.status.ok:
            error = resp.status.err
            print('创建 session 失败', error)
        return resp.data

    def close_session(self, session: Session) -> None:
        req = service_pb2.SessionFinishRequest(
            id=session.id,
            date_end=int(datetime.now().timestamp())
        )
        resp = self.stub.finish_session(req)

        if not resp.status.ok:
            print("关闭会话失败: ", resp.status.err)

    # def wait_for_kill_session_message(self):
    #     def dispatch_task():
    #         try:
    #             response_iterator = self.stub.dispatchTask(iter([]))
    #             for task_response in response_iterator:
    #                 target_session = None
    #                 for session in SessionManager.get_store().values():
    #                     if isinstance(session, JMSSession):
    #                         if session.get_jms_session().get_id() == task_response.task.session_id:
    #                             target_session = session
    #                             break
    #                 if target_session is not None:
    #                     if task_response.task.action == KillSession:
    #                         target_session.close()
    #                     req = self.stub.FinishedTaskRequest(task_id=task_response.task.id)
    #                     request_observer.on_next(req)
    #         except Exception as e:
    #             print("waitSessionMessage error: {}".format(e))
    #         else:
    #             print("waitSessionMessage completed")
    #
    #     request_observer = self.stub.dispatchTask(iter([]))
    #     thread = threading.Thread(target=dispatch_task)
    #     thread.start()


class JMSSession(BaseWisp):
    def __init__(self, session: Session, token_resp: TokenAuthInfo):
        super().__init__()
        self.jms_session = session
        self.command_acls = list(token_resp.filter_rules)
        self.expire_time = token_resp.expire_info.expire_at
        self.max_idle_time_delta = token_resp.setting.max_idle_time
        self.session_handler = None
        self.command_handler = None
        self.replay_handler = None

    def get_username(self) -> str:
        return self.jms_session.user

    def get_datasource_name(self) -> str:
        return self.jms_session.asset

    def active_session(self) -> None:
        self.session_handler = SessionHandler()
        self.command_handler = CommandHandler(self.jms_session, self.command_acls)
        self.replay_handler = ReplayHandler(self.jms_session)

    def close(self) -> None:
        self.replay_handler.upload()
        self.session_handler.close_session(self.jms_session)

    # TODO chat_func 还没有写
    def with_audit(self, command: str, chat_func):
        command_record = CommandRecord(input=command)
        try:
            self.command_handler.command_acl_filter(command_record)
            self.replay_handler.write_input(command_record.input)

            result = chat_func.run()
            command_record.output = result

            self.replay_handler.write_output(result.output)
            return result

        except Exception as e:
            command_record.error = str(e)
            self.replay_handler.write_output(str(e))
            raise e

        finally:
            self.command_handler.record_command(command_record)
