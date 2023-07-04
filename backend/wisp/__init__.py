import os
import sys

import grpc

from api import globals


def setup_protobuf():
    current_dir = os.path.dirname(os.path.abspath(__file__))
    protobuf_path = os.path.join(current_dir, 'protobuf')
    sys.path.insert(0, protobuf_path)

    globals.GRPC_CHANNEL = grpc.insecure_channel('localhost:9090')


def shutdown_protobuf():
    if globals.GRPC_CHANNEL is None:
        return
    globals.GRPC_CHANNEL.close()
