from wisp.protobuf.common_pb2 import Session


class SessionManager:
    store = {}

    @staticmethod
    def register_session(session: Session):
        session_id = f'{session.id}'
        SessionManager.store[session_id] = session
        return session_id

    @staticmethod
    def unregister_session(session: Session):
        if session.id in SessionManager.store:
            SessionManager.store.pop(session.id, None)

    @staticmethod
    def get_store():
        return SessionManager.store

    @staticmethod
    def get_session(session_id):
        return SessionManager.store.get(session_id)
