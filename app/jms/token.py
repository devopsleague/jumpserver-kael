from wisp.protobuf import service_pb2
from wisp.exceptions import WispError
from wisp.protobuf.common_pb2 import TokenAuthInfo

from jms.base import BaseWisp
from utils.logger import get_logger

logger = get_logger(__name__)


class TokenHandler(BaseWisp):

    def get_token_auth_info(self, token: str) -> TokenAuthInfo:
        req = service_pb2.TokenRequest(token=token)
        token_resp = self.stub.GetTokenAuthInfo(req)
        if not token_resp.status.ok:
            error_message = f'Failed to get token: {token_resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)
        return token_resp.data
