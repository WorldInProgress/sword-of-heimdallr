# Messaging Spec

本文档介绍了通信基座的基本通信机制以及服务间、进程间的消息传递规范。gRPC 协议和 ZeroMQ 库提供了消息传递的传输底层。

## Versioning

当前版本为**0.1**，即体系构建阶段的临时版本。

更新内容为：

- 初始构建

## Introduction

此处为本项目微服务架构简介，应涉及服务间的交互机制。

## General Message Format

一条消息的通用格式应由以下 5 个字典要素组合而成。

### Header

```json
{
    "msg_id": str,            # UUID, must be unique per message
    "timestamp": str,         # ISO 8601 timestamp for when the message is created
    "msg_type": str,          # All recognized message type strings are listed below
    "compression": enum,      # REQUEST || RESPONSE || ACK
    "encoding": enum,         # gzip || snappy || none
    "version": "0.1",         # the message protocol version
}
```

### Meta

```json
{
    "service": str,       # UUID, optional, only when this message is a response
    "method": str,        # UUID, optional, only when this message is a response
    "parent_msg_id": str, # UUID, optional, only when this message is a response
    "node_id": str,       # UUID, must be unique per node
    "user_id": str,       # UUID, must be unique per user
    "priority": enum,     # HIGH || MEDIUM || LOW
}
```

### Content

```json
{}
```

### Security

```json
{
    "msg_id": str,        # UUID, must be unique per message
    "parent_msg_id": str, # UUID, optional, only when this message is a response
    "node_id": str,       # UUID, must be unique per node
    "user_id": str,       # UUID, must be unique per user
    "timestamp": str,     # ISO 8601 timestamp for when the message is created
    "msg_type": str,      # All recognized message type strings are listed below
    "compression": str,   # gzip || snappy || none
    "encoding": str,      # utf-8 || ascii || none
    "version": "0.1",     # the message protocol version
}
```

## The Wire Protocol

```json
[
    b"u-u-i-d",           # zmq identity(ies)
    b"<IDS|MSG>",         # delimiter
    b"baddad42",          # HMAC signature
    b"{header}",          # serialized header dict
    b"{parent_header}",   # serialized parent header dict
    b"{metadata}",        # serialized metadata dict
    b"{content}",         # serialized content dict
    b"\xf0\x9f\x90\xb1"   # extra raw data buffer(s)
    # ...
]
```

delimiter 是分隔符

分隔符之前是 zmq 的路由前缀，可用作消息的 topic，大多数情况下会被 msg_type 取代功能

## Message Mode

![image-20250225100506906](C:\Users\DMK\AppData\Roaming\Typora\typora-user-images\image-20250225100506906.png)

在 **ZeroMQ（ØMQ）** 中，`ROUTER` 和 `DEALER` 是 **高级套接字模式**，用于构建 **灵活的异步消息路由**。它们主要用于 **扩展 Request-Reply（请求-响应）模式**，支持多对多、异步处理等复杂通信需求。

| **Socket 类型** | **作用**               | **特点**                                                                               |
| --------------- | ---------------------- | -------------------------------------------------------------------------------------- |
| `ROUTER`        | 高级响应端（如服务器） | 允许多个客户端连接，维护客户端身份标识（identity），可任意路由消息，并主动选择响应对象 |
| `DEALER`        | 高级请求端（如客户端） | 支持多对多通信，异步处理，可负载均衡请求                                               |

**ROUTER/DEALER 组合：**

- `ROUTER`（路由器）可以管理多个 `DEALER`（处理者）。
- `DEALER` 可以在多个 `ROUTER` 之间分配负载，实现并行任务处理。

## Messages on the shell (ROUTER/DEALER) channel

### Execute

#### Execute request——Message type: `execute_request`

```json
content = {
  'code' : str,
  'silent' : False,
  'store_history' : True,
  'user_expressions' : dict,
  'allow_stdin' : True,
  'stop_on_error' : True,
}
```

```json
'code'   # Source code to be executed by the kernel, one or more lines.
'silent'  # if True, no output or execute_result, 'store_history' to be     False, has no effect on 'user_expressions'
'store_history'  # Whether to save in execution history
'user_expressions' # extra expressions to be executed, and return the results
'allow_stdin'  # Whether to allow the kernel to request user input
'stop_on_error'  # Whether to abort the execution queue when error encountered
```

```json
user_expressions_result = {
    [expression_name]: {
        "status": "ok" || "error",
        "data": {[type]: [value]}
    }
}
```

#### Execution results——Message type: `execute_reply`

```json
content = {
  'status' : 'ok',
  'execution_count' : int,
  'payload' : list(dict),
  'user_expressions' : dict,
}
```

```json
content = {
  'status' : 'error',
  'execution_count' : int,
}
```

```json
'status'  # One of: 'ok' OR 'error'
'execution_count' # The global kernel counter that increases by one with each request    that stores history.  This will typically be used by clients to    display prompt numbers to the user.  If the request did not store    history, this will be the current value of the counter in the     kernel.
'payload'  # a way to trigger frontend actions
{
    "source": "page",  # 分页显示长文本
    # mime-bundle of data to display in the pager.
    # Must include text/plain.
    "data": mimebundle,
    # line offset to start from
    "start": int,
}
{
    "source": "set_next_input", # 创建新单元格或在命令行中创建新的输入
    # the text contents of the cell to create
    "text": "some cell content",
    # If true, replace the current cell in document UIs instead of inserting
    # a cell. Ignored in console UIs.
    "replace": bool,
}
{
    "source": "edit_magic",  # 打开文件进行编辑
    "filename": "/path/to/file.py",  # the file to edit
    "line_number": int,  # the line number to start with
}
{
    "source": "ask_exit",  # 提示用户退出
    # whether the kernel should be left running, only closing the client
    "keepkernel": bool,
}
```

### History

#### Message type: `history_request`

#### Message type: `history_reply`

### Comm Info

Jupyter 的 **Comm** 机制允许前端（如 Jupyter Notebook、JupyterLab、IPython 控制台）和内核之间进行**双向通信**，而不局限于代码执行。

- **ipywidgets**（Jupyter 小部件）：前端 UI 控件（如滑块、按钮）和 Python 代码交互。

- **自定义前后端通信**：前端 JavaScript 和 Python 代码实时传输数据。

- **动态数据更新**：如绘图库（Matplotlib、Plotly）可以使用 Comm 实时更新图表。

#### Message type: `comm_info_request`

#### Message type: `comm_info_reply`

### Kernel Info

#### Message type: `kernel_info_request`

#### Message type: `kernel_info_reply`

## Messages on the Control (ROUTER/DEALER) channel

### Kernel Shutdown

关闭

### Kernel Interrupt

中断当前执行

### Debug

调试

## Messages on the IOPub (PUB/SUB) channel

### Streams (stdout, stderr, etc)——Message type: `stream`

```json
content = {
    # The name of the stream is one of 'stdout', 'stderr'
    'name' : str,

    # The text is an arbitrary string to be written to that stream
    'text' : str,
}
```

### Display Data——Message type: `display_data`

```json
content = {

    # The data dict contains key/value pairs, where the keys are MIME
    # types and the values are the raw data of the representation in that
    # format.
    'data' : dict,

    # Any metadata that describes the data
    'metadata' : dict,

    # Information not to be persisted to a notebook or other documents. Intended to  # live only during a live kernel session.
    'transient': dict,
}
```

MIME

```json
{
  "data": {
    "text/plain": "Hello, Jupyter!",
    "text/html": "<b>Hello, Jupyter!</b>",
    "image/png": "iVBORw0KGgoAAAANSUhEUgAA...",
    "application/json": "{'key': 'value'}"
  },
  "metadata": {
    "image/png": {
      "width": 640,
      "height": 480
    }
  }
}
```

### Update Display Data——Message type: `update_display_data`

更新数据展示信息，格式同 display_data，其中 transient 指定 display_id 变为必须字段

### Code inputs——Message type: `execute_input`

获取执行过程的源代码

### Execution results——Message type: `execute_result`

```json
content = {
    'execution_count' : int,
    'data' : dict,
    'metadata' : dict,
}
```

data 与 metadata 和 display_data 结构一致

### Execution errors——Message type: `error`

```json
content = {
   # Similar content to the execute_reply messages for the 'error' case,
   # except the 'status' and 'execution_count' fields are omitted.
}
```

### Kernel status——Message type: `status`

```json
content = {
    # When the kernel starts to handle a message, it will enter the 'busy'
    # state and when it finishes, it will enter the 'idle' state.
    # The kernel will publish state 'starting' exactly once at process startup.
    'execution_state' : ('busy', 'idle', 'starting')
}
```

## Messages on the stdin (ROUTER/DEALER) channel

kernel 根据代码发送请求，要求前端提示用户输入

### Message type: `input_request`

```json
content = {
    # the text to show at the prompt
    'prompt' : str,
    # Is the request for a password?
    # If so, the frontend shouldn't echo input.
    'password' : bool
}
```

### Message type: `input_reply`

```json
content = { 'value' : str }
```

## Custom Messages

自定义消息，通过引入`Comm`，在前端和 kernel 中都有，实现双向通信

这些消息是完全对称的 - 内核和前端都可以发送每条消息，并且没有消息需要回复。内核在 Shell 通道上监听这些消息，而前端在 IOPub 通道上监听这些消息。

### Opening a Comm——Message type: `comm_open`

```json
{
  "comm_id": "u-u-i-d",
  "target_name": "my_comm",
  "data": {}
}
```

前段请求针对 target_name 建立 comm，后端分配 comm_id，

### Comm Messages——Message type: `comm_msg`

```json
{
  "comm_id": "u-u-i-d",
  "data": {}
}
```

### Tearing Down Comms——Message type: `comm_close`

```json
{
  "comm_id": "u-u-i-d",
  "data": {}
}
```

### Output Side Effects

- 响应信息应设置 parent header 以便追溯
- 处理请求应设置 busy / idle 状态

## Changelog

### 0.1

- 学习 Messaging in Jupyter
- 初始化规范
