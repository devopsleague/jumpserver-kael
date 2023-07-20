from api import globals
from wisp.protobuf import service_pb2_grpc


class BaseWisp:

    def __init__(self):
        grpc_channel = globals.grpc_channel
        self.stub = service_pb2_grpc.ServiceStub(grpc_channel)
