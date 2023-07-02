from asgiref.typing import ASGI3Application


class AuthMiddleware:
    def __init__(self, app: ASGI3Application, my_option: str):
        self.app = app
        self.my_option = my_option

    async def __call__(self, scope, receive, send):
        print(f"MyMiddleware initialized with option: {self.my_option}")

        await self.app(scope, receive, send)
