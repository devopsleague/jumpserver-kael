from enum import Enum, auto


class WSStatusCode(Enum):
    normal_close = 1000             # 正常关闭连接。表示连接已成功关闭。
    terminal_leave = auto()         # 终端离开。表示终端（客户端）离开或导致连接关闭。
    protocol_error = auto()         # 协议错误。表示由于协议错误而导致连接关闭。
    unacceptable_data = auto()      # 不可接受的数据。表示接收到不可接受的数据类型。
    no_status_received = 1005       # 空状态。表示未收到任何状态码。
    connection_closed = auto()      # 连接关闭。表示连接意外关闭，原因不明确。
    data_error = auto()             # 数据错误。表示收到的数据格式不正确。
    message_too_large = auto()      # 消息过大。表示收到的消息过大而无法处理。
    too_many_messages = auto()      # 接收到的消息过多。表示收到的消息队列超出处理能力。
    extension_error = auto()        # 扩展错误。表示由于收到的扩展不符合预期而导致连接关闭。
    server_error = auto()           # 服务器错误。表示由于服务器遇到不可预知的情况而关闭连接。
    # 1012-2999 保留用于未来定义的状态码。
    # 3000-3999保留用于应用程序定义的状态码。
