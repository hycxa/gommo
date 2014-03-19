# 使用golang设计的分布式框架

## 目标

- 无限的在线人数
- 每玩家每秒10次普通响应
- 每玩家每秒1次高压力响应
- 每玩家每秒1次玩家间互动
- 每玩家每秒1次社交团体间互动
- 无限的战斗场景
- 无限的AI数量
- 每秒一次重度AI
- 每秒5次轻度AI

## 设计思路

参考erlang语言并发模型，整个框架分为以下几层：
server
routine
state


### Routine

routine之间

### Message

#### P2P

#### Broadcast

- 所有现实中的对象都对应程序中一个Object
- 所有Object之间都通过Message传递。
- 所有的Message的传递、响应都是异步的。
- Message分为P2P和Broadcast两种。
- P2PMessage
    + Sender必须拥有Receiver的Name
    + Sender和Receiver中间可以定制Filter
    + Receiver接收到的任何消息必须返回一个Response（系统调用Reply）
    + Sender的每一次Send必须处理ReceiverNotFound消息
    + Sender的每一次Send必须处理Timeout消息
    + Sender的每一次Send必须处理Response消息
    + 消息发送流程如下：
        * Sender调用Send接口，设置ReceiverNotFoundFunction，TimeoutFunction，ResponseFunction
        * 系统未找到Receiver则向Sender发送ReceiverNotFound，流程结束
        * Receiver收到任何消息返回对应消息的Response
        * Sender处理Response，流程结束
        * 从Sender调用Send的时刻起间隔Timeout时间，系统向Sender发送Timeout消息，流程结束
        * 以上流程意味着Response在Timeout后收到不会被响应，因此Response不能用于处理逻辑，只能用于确定消息传递正常完成。（这样似乎通过Timeout就行了，可能不需要把Response发给Sender，甚至Response可以由系统自动产生）
- Broadcast可以指定任意数量的Object，同时不需要处理Timeout、Response、NotFound
- 所有Object本身都是一个FSM，根据当前State和收到的消息做状态切换
- FSM是一颗自包含的树
