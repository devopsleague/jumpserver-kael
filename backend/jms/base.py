import asyncio


class BaseHandler:
    async def run(self, method: str, *args, **kwargs):
        if not hasattr(self, method):
            return
        func = getattr(self, method)
        asyncio.create_task(func(*args, **kwargs))
