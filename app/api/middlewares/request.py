from asgiref.typing import ASGI3Application


class RequestMiddleware:
    def __init__(self, app: ASGI3Application):
        self.app = app

    async def __call__(self, scope, receive, send):
        await self.app(scope, receive, send)
