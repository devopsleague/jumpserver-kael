import asyncio

from api import globals
from wisp.protobuf import service_pb2_grpc


class BaseWisp:

    def __init__(self):
        grpc_channel = globals.grpc_channel
        self.stub = service_pb2_grpc.ServiceStub(grpc_channel)

    async def run(self, method: str, *args, **kwargs):
        if not hasattr(self, method):
            return
        func = getattr(self, method)
        asyncio.create_task(func(*args, **kwargs))

    def sync_run(self, method: str, *args, **kwargs):
        if not hasattr(self, method):
            return
        func = getattr(self, method)
        return func(*args, **kwargs)
