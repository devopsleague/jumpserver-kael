from .handler import JMSSession


class SessionManager:
    store = {}

    @staticmethod
    def register_jms_session(jms_session: JMSSession):
        session_id = jms_session.session.id
        SessionManager.store[session_id] = jms_session
        return session_id

    @staticmethod
    def unregister_jms_session(jms_session: JMSSession):
        if jms_session.session.id in SessionManager.store:
            SessionManager.store.pop(jms_session.session.id, None)

    @staticmethod
    def get_store():
        return SessionManager.store

    @staticmethod
    def get_jms_session(session_id: str):
        return SessionManager.store.get(session_id)
