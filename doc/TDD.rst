====================
golang分布式服务器框架
====================

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
-------

参考erlang语言并发模型，基于消息机制。

整个框架包含以下几个层次，包含关系为Server->Routine->FSM：

Server
======

一个Server对应一个操作系统进程，主要作用为从网络层接收消息并启动Routine，包括以下几类：

Gate
	和Client直接连接，每个Client一个Routine接收消息，过滤，转发给App服务器。Gate上还有若干种Org用于消息转发或广播。

App
	运行应用逻辑的服务器，依然是每个Client对应一个Routine，进行消息的处理。

Routine
=======

一个Routine对应一个Processer，有以下几类：

Person
	直接对应一个Client，所有的合法消息先到这里来处理，然后再回复或者转发给其他Processer。

Org
	对应一组Client，用于广播消息，处理组织类消息（例如公会）。

FSM
===

FSM是一颗自包含的树，一个Routine中包含一颗FSM Tree分发、处理消息。这里已经涉及到具体的应用逻辑。

Message
=======

+ 所有的Message的传递、响应都是异步的。
+ Sender和Receiver可以在同一进程也可在不同进程。
+ Sender必须知道Receiver的ID
+ Sender和Receiver中间可以定制Filter（？）
+ 任何一个消息的定义必须有四个结构，例如：

.. code:: golang

	// Message
	type MLogin struct {
		ID	UID
	}

	// Response
	type RLogin bool

	// ReceiverNotFound
	type NLogin UID

	// Timeout
	type TLogin int64

+ Receiver接收到的任何消息必须返回一个Response
+ Sender的每一次Send必须处理ReceiverNotFound消息（？）
+ Sender的每一次Send必须处理Timeout消息（？）
+ Sender的每一次Send必须处理Response消息
+ 消息发送流程如下：
	* Sender调用Notify接口：
	.. code:: golang

		func Notify(sender Notifier, receiver Processer, m Message) (ok bool) {
			return true
		}

	* 系统未找到Receiver则向Sender发送ReceiverNotFound，流程结束
	* Receiver收到任何消息返回对应消息的Response
	* Sender处理Response，流程结束
	* 从Sender调用Send的时刻起间隔Timeout时间，系统向Sender发送Timeout消息，流程结束
	* 以上流程意味着Response在Timeout后收到不会被响应，因此Response不能用于处理逻辑，只能用于确定消息传递正常完成。（这样似乎通过Timeout就行了，可能不需要把Response发给Sender，甚至Response可以由系统自动产生）

实现细节
-------

* 消息传递使用chan
