from http import cookies
from asgiref.typing import ASGI3Application

from api import globals


class I18nMiddleware:
    def __init__(self, app: ASGI3Application):
        self.app = app

    @staticmethod
    def parse_cookie(scope):
        headers = dict(scope.get('headers'))
        cookie_str = headers.get(b'cookie', b'').decode()

        cookie = cookies.SimpleCookie()
        cookie.load(cookie_str)
        return {key: morsel.value for key, morsel in cookie.items()}

    async def __call__(self, scope, receive, send):
        try:
            cookie = self.parse_cookie(scope)
            globals.language = cookie.get('django_language', 'zh-hans')
        except Exception:
            pass
        await self.app(scope, receive, send)
