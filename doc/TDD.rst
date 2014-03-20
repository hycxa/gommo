======================
golang分布式服务器框架
======================

目标
----

- 单一大服
- 无限的在线人数
- 每玩家每秒10次普通响应
- 每玩家每秒1次高压力响应
- 每玩家每秒1次玩家间互动
- 每玩家每秒1次社交团体间互动
- 无限的战斗副本
- 无限的AI数量
- 每秒一次重度AI
- 每秒5次轻度AI

设计思路
--------

参考erlang语言并发模型，基于消息机制。

整个框架包含以下几个层次，包含关系为Node->Process->FSM：

Node
====

一个Node对应一个操作系统进程，主要作用为从网络层接收消息并启动Process，包括以下几类：

Gate
	和Client直接连接，每个Client一个GateProcess接收消息，过滤，转发给对应的AppProcess服务器。

App
	运行应用逻辑的服务器，依然是每个Client对应一个Process，进行消息的处理。

Process
=======

一个Node包含多个Process，有以下几类：

Person
	直接对应一个Client，所有的合法消息先到这里来处理，然后再回复或者转发给其他Node。

Org
	对应一组Client，用于广播消息，处理组织类消息（例如公会、队伍、战斗）。

FSM
===

FSM是一颗自包含的树，一个Process中包含一颗FSM Tree分发、处理消息。这里已经涉及到具体的应用逻辑。

Message
=======

+ 所有的Message的传递、响应都是异步的。
+ Sender和Receiver都是Process。
+ Sender和Receiver可以在同一Node也可在不同Node。
+ Sender必须知道Receiver的ID才能发消息。
+ 消息的定义可以有Timeout，例如：

.. code:: go

	// Message
	type MLogin struct {
		ID	UID
	}

	type Timeout interface {
		time.Duration Interval()
	}

	// Timeout
	type TLogin Timeout

+ Sender可以选择是否处理Timeout
+ 消息发送流程如下：
	* Sender调用Notify接口：

	.. code:: go

		func Notify(sender UUID, receiver UUID, message interface{}, t Timeout)

	* 从Sender调用Notify的时刻起间隔timeout时间，Sender收到Timeout消息，流程结束

实现细节
--------

* 消息传递使用chan
* 
