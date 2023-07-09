import re
import time
from typing import List, Optional
from datetime import datetime
from pydantic import BaseModel
from starlette.websockets import WebSocket

from jms.base import BaseWisp
from api.schemas import AskResponse, AskResponseType
from wisp.protobuf import service_pb2
from wisp.exceptions import WispError
from wisp.protobuf.common_pb2 import Session, CommandACL, RiskLevel
from utils import reply
from utils.logger import get_logger

logger = get_logger(__name__)


class CommandRecord(BaseModel):
    input: Optional[str] = None
    output: Optional[str] = None
    risk_level: str = RiskLevel.Normal


class CommandHandler(BaseWisp):
    REJECT_MESSAGE = "reject by acl rule"
    WAIT_TICKET_TIMEOUT = 60 * 3
    WAIT_TICKET_INTERVAL = 3

    def __init__(self, websocket: WebSocket, session: Session, command_acls: List[CommandACL]):
        super().__init__()
        self.session = session
        self.websocket = websocket
        self.command_acls = command_acls

    async def record_command(self, command_record: CommandRecord):
        req = service_pb2.CommandRequest(
            sid=self.session.id,
            org_id=self.session.org_id,
            asset=self.session.asset,
            account=self.session.account,
            user=self.session.user,
            timestamp=int(datetime.timestamp(datetime.now())),
            input=command_record.input,
            output=command_record.output,
            risk_level=command_record.risk_level
        )
        resp = self.stub.UploadCommand(req)
        if not resp.status.ok:
            error_message = f'Failed to upload command: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)

    def match_rule(self, command: str):
        for command_acl in self.command_acls:
            for command_group in command_acl.command_groups:
                flags = re.UNICODE
                if command_group.ignore_case:
                    flags |= re.IGNORECASE
                try:
                    pattern = re.compile(command_group.pattern, flags)
                    if pattern.search(command.lower()) is not None:
                        return command_acl
                except re.error as e:
                    error_message = f'Failed to re invalid pattern: {command_group.pattern} {e}'
                    logger.error(error_message)
                    # TODO Exception
                    raise Exception(error_message)

    async def create_and_wait_ticket(self, command: str, command_acl: CommandACL) -> bool:
        req = service_pb2.CommandConfirmRequest(
            cmd=command,
            session_id=self.session.id,
            cmd_acl_id=command_acl.id
        )
        resp = self.stub.CreateCommandTicket(req)
        if not resp.status.ok:
            error_message = f'Failed to create ticket: {resp.status.err}'
            logger.error(error_message)
            raise WispError(error_message)

        return await self.wait_for_ticket_status_change(resp.info)

    async def wait_for_ticket_status_change(self, ticket_info: service_pb2.TicketInfo):
        await reply(
            self.websocket, AskResponse(
                type=AskResponseType.waiting,
                conversation_id=self.session.id,
                system_message=f'等待工单审批: {ticket_info.ticket_detail_url}'
            )
        )
        start_time = time.time()
        end_time = start_time + self.WAIT_TICKET_TIMEOUT

        ticket_closed = False
        is_continue = False
        while time.time() <= end_time:
            check_request = service_pb2.TicketRequest(req=ticket_info.check_req)
            check_response: service_pb2.TicketStateResponse = self.stub.CheckTicketState(check_request)

            if not check_response.status.ok:
                error_message = f'Failed to check ticket status: {check_response.status.err}'
                logger.error(error_message)
                break

            state = check_response.Data.state
            if state == service_pb2.TicketState.Approved:
                is_continue = True
                ticket_closed = True
                break
            elif state in [service_pb2.TicketState.Rejected, service_pb2.TicketState.Closed]:
                ticket_closed = True
                await reply(
                    self.websocket, AskResponse(
                        type=AskResponseType.waiting,
                        conversation_id=self.session.id,
                        system_message=f'工单关闭或拒绝'
                    )
                )
                break

            time.sleep(self.WAIT_TICKET_INTERVAL)

        if not ticket_closed:
            self.close_ticket(ticket_info)

        return is_continue

    async def command_acl_filter(self, command: CommandRecord):
        is_continue = False
        acl = self.match_rule(command.input)
        if acl is not None:
            command.risk_level = RiskLevel.Danger
            if acl.action == CommandACL.Reject:
                await reply(
                    self.websocket, AskResponse(
                        type=AskResponseType.reject,
                        conversation_id=self.session.id,
                        system_message=self.REJECT_MESSAGE
                    )
                )
            elif acl.action == CommandACL.Review:
                try:
                    is_continue = self.create_and_wait_ticket(command.input, acl)
                except Exception as e:
                    print(command.input, str(e))
            else:
                is_continue = True
        return is_continue

    def close_ticket(self, ticket_info: service_pb2.TicketInfo):
        req = service_pb2.TicketRequest(req=ticket_info.cancel_req)
        resp = self.stub.CancelTicket(req)
        if not resp.status.ok:
            error_message = f'Failed to close ticket: {resp.status.err}'
            logger.error(error_message)
