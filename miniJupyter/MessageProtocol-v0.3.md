# Message Protocol

本文档介绍了通信基座的基本通信机制以及服务间、进程间的消息传递规范。gRPC 协议和 ZeroMQ 库提供了消息传递的传输底层。

## Version

当前版本为**0.3**，即初版可使用规范基础上进行的性能升级。

## Introduction

此处为本项目后端通信协议介绍，涉及进程间、服务间的交互机制。

## General Message Format

一条消息的通用格式应由以下 5 个字典要素组合而成。

```json
{
    "header" : {
        "msg_id": "...",
        "msg_type": "...",
        ...
    },
    "parent_header": {},
    "meta": {},
    "content": {},
    "security": {},
}
```

### Header

负责协议级别的控制信息，确保消息的正确传输和解析

```json
{
    "msg_id": str,            # UUID, must be unique per message
    "session_id": str,    	  # UUID, must be unique per session
    "user_id": str,           # UUID, must be unique per user
    "timestamp": str,         # ISO 8601 timestamp for when the message is created
    "msg_type": str,          # All recognized message type strings are listed below
    "compression": enum,      # gzip || snappy || none
    "encoding": enum,         # json || protobuf || custom
    "transport": enum,		  # zmq || gRPC
    "version": "0.3",         # the message protocol version
}
```

### Parent Header

当消息是响应时，记录请求消息的Header

```json
{
    # parent_header is a copy of the request's header
    'msg_id': '...',
    ...
}
```

### Meta

存储业务层面的附加信息，方便解析和扩展

```json
{
    "priority": enum,     # HIGH || NORMAL || LOW
    "tags": list,         # extra optional info
}
```

### Content

主要存储实际的业务数据，具体结构由msg_type决定

```json
{}
```

### Security

```json
{
    "token": str,         # JWT token
    "encryption": enum,   # AES || RSA || None
}
```

## Message Type

### ROUTER / DEALER

#### Execute

##### `execute_request`

```json
content = {
    "command_id": str,       	# UUID, must be unique per command
    "service": str,         	# service which contain the needed method
    "method": str,          	# specific method to call
    "params": {             	# parameters key-values
        "key": "value",
    },
	"condition": {				# condition for command to execute
		"key": "value"
	},
    "dependency": [str],		# commands to depend on
    "timeout": num,				# Task timeout(ms)
    "retry": {
        "max_attempts": num,	# max retry attempts
        "strategy": enum,		# exponential_backoff
    },
	"stop_on_error": bool,		# whether to stop when error encountered
	"allowed_users": [str],		# users allowed to suscribe this command result
}
```

##### `execute_reply`

执行一条命令，可能有参数问题立刻报错error，若没有依赖则立刻starting，若有依赖的命令则进入waiting

```json
content = {
  "status": enum,		# error || starting || waiting
}
```

#### Query

##### `core_info_request`

```json
content = {}
```

##### `core_info_reply`

```json
content = {
    "status": enum,			# ok || error
    "core_status": enum,	# healthy || down
    "core_version": str,
    "cpu_usage": str,
    "memory_usage": str,
    "disk_usage": str,
    "network_io": str,
    "active_connections": num,
    "running_tasks": num,
    "task_queue_size": num,
}
```

### XPUB / XSUB + PUB / SUB

XPUB/XSUB是PUB/SUB的消息中介，可支持订阅者权限控制，仅允许有权限的用户订阅特定主题

| 角色 | 作用                                    |
| ---- | --------------------------------------- |
| PUB  | 负责**发送消息**，不关心谁会接收        |
| XPUB | **记录订阅**，决定消息是否转发给订阅者  |
| XSUB | **中继订阅请求**，将订阅信息传递给 XPUB |
| SUB  | 订阅消息，接收符合订阅条件的消息        |



| **角色 (Role)** | 允许订阅的主题 (Allowed Topics)   |
| --------------- | --------------------------------- |
| admin           | `system.*`, ``user.*`, `public.*` |
| moderator       | `user.*`, `public.*`              |
| user            | `public.*`                        |
| guest           | `public.announcements`            |

#### Result

##### `execute_result`

```json
content = {
    "status": enum,		# success || error
    "result": {},		# result data to show or error info 
}
```

##### `stream`

```json
content = {
    "type": enum,		# stdout || stderr
    "text": str,		# arbitrary string to be written to that stream
}
```

## Heartbeat

心跳**不遵循**通用消息格式，仅需简单的字符串通信，分为双向心跳监测

### PUB / SUB

##### `core_heart_beat`

```json
"text": "core alive"
```

##### `client_heart_beat`

```json
"text": "[session_id] alive"
```

## The Wire Protocol

```json
[
    b"u-u-i-d",           # zmq identity(ies)
    b"<IDS|MSG>",         # delimiter
    b"baddad42",          # HMAC signature
    b"{header}",          # serialized header dict
    b"{parent_header}",   # serialized parent header dict
    b"{meta}",        	  # serialized meta dict
    b"{content}",         # serialized content dict
    # ...
]
```

delimiter 是分隔符

分隔符之前是 zmq 的路由前缀，可用作消息的 topic

## Custom Messages

自定义消息，通过引入`Comm`，在前端和 kernel 中都有，实现双向通信

这些消息是完全对称的 - 内核和前端都可以发送每条消息，并且没有消息需要回复。内核在 Shell 通道上监听这些消息，而前端在 IOPub 通道上监听这些消息。

### `comm_open`

```json
{
  "comm_id": "u-u-i-d",
  "target_name": "my_comm",
  "data": {}
}
```

前段请求针对 target_name 建立 comm ，后端分配 comm_id ，多个 comm 可以共享一个target_name

### `comm_msg` || `comm_close`

```json
{
  "comm_id": "u-u-i-d",
  "data": {}
}
```

响应消息应添加parent_header字段，便于追溯，模仿通用消息

```json
{
  "parent_header": {
      "comm_id": "u-u-i-d",
  	  "data": {}
  },
  "comm_id": "u-u-i-d",
  "data": {}
}
```

## Changelog

### 0.1

- 学习 Messaging in Jupyter
- 初始化规范

### 0.2

- 针对项目适应性调整
- 初版可使用规范

### 0.3

- 调整了execute_request部分，变为单条有依赖命令