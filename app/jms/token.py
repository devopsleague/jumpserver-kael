from wisp.protobuf import service_pb2
from wisp.protobuf.common_pb2 import TokenAuthInfo

from jms.base import BaseWisp


class TokenHandler(BaseWisp):

    def get_token_auth_info(self, token: str) -> TokenAuthInfo:
        req = service_pb2.TokenRequest(token=token)
        token_resp = self.stub.GetTokenAuthInfo(req)
        if not token_resp.status.ok:
            error = token_resp.status.err
            print('获取 token 失败', error)

        return token_resp.data
