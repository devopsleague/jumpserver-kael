import os
import sys

import grpc

from api import globals


def setup_protobuf():
    current_dir = os.path.dirname(os.path.abspath(__file__))
    protobuf_path = os.path.join(current_dir, 'protobuf')
    sys.path.insert(0, protobuf_path)

    globals.grpc_channel = grpc.insecure_channel('localhost:9090')


def shutdown_protobuf():
    if globals.grpc_channel is None:
        return
    globals.grpc_channel.close()


setup_protobuf()
