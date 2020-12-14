from .cybervein_pb2 import (
    AuthRequest,
    Token,
    CommandRequest,
    QueryResponse,
    QueryPrivateWithAddrRequest,
    ExecuteResponse,
    ExecuteAsyncResponse
)

from .cybervein_pb2_grpc import (
    cyberveinStub
)

import grpc


def get_client(server):
    conn = grpc.insecure_channel(server)
    return cyberveinStub(channel=conn)
