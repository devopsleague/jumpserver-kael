from starlette import status
from starlette.responses import Response
from fastapi import HTTPException, APIRouter
from api.jms import JMSSession, SessionManager
from api.schemas import Conversation, JMSState
from utils.logger import get_logger

logger = get_logger(__name__)
router = APIRouter()


def get_jms_session(session_id: str) -> JMSSession:
    jms_session = SessionManager.get_jms_session(session_id)
    if jms_session and isinstance(jms_session, JMSSession):
        return jms_session
    else:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail='Not found conversation')


@router.post('/interrupt_current_ask/')
async def interrupt_current_ask(conversation: Conversation):
    jms_session = get_jms_session(conversation.id)
    jms_session.current_ask_interrupt = True
    return Response(status_code=status.HTTP_200_OK)


@router.post('/jms_state/')
async def jms_state(state: JMSState):
    jms_session = get_jms_session(state.id)
    jms_session.jms_state.activate_review = state.activate_review
    return Response(status_code=status.HTTP_200_OK)
