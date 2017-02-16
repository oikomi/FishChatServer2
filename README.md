# FishChatServer2

[![Build Status](https://travis-ci.org/oikomi/FishChatServer2.svg?branch=master)](https://travis-ci.org/oikomi/FishChatServer2)
[![Go Report Card](https://goreportcard.com/badge/github.com/oikomi/FishChatServer2)](https://goreportcard.com/report/github.com/oikomi/FishChatServer2)

目录
=================
* [1说明](#1说明)
* [2特性](#2特性)
* [3架构](#3架构)
* [4协议](#4协议)
* [5数据模型](#5数据模型)
* [6服务说明](#6服务说明)
* [7依赖](#7依赖)
* [8部署](#8部署)
    * [8.1普通部署](#8.1普通部署)
    * [8.2容器部署](#8.2容器部署)
* [9测试](#9测试)
* [10监控](#10监控)


1说明
======
吸取了第一版的经验以及我们商业版的探索. 第二版我更多的思考在不要过多的实现轮子, 这个版本将很多业务无关的代码去掉, 并且尽量靠拢`微服务`.
部署方式可以支持：
> * 普通部署
> * 容器部署 (Kubernetes + Docker)


2特性
======


**[⬆ 回到顶部](#目录)**

3架构
======

![](./doc/architecture.png)

### 3.1聊天设计方案

这里简单陈述一下消息的设计思路, 如下

![](./doc/msg.png)

其中:

1 消息是通过版本号维护的, 是一个递增的int64整数, 由`idgen`提供服务, 每个用户都有一个独立的自增id

2 新消息到达, 服务端只负责发送给客户端一个轻量级的notify通知

3 客户端收到notify后, 发起同步请求


### 3.2存储方案

其中最关键的是HBase存储, 所有的消息通过Kafka消费后将插入HBase中, 消息存储的时候会带上递增的版本号. 这样客户端带上一个版本号来请求消息的时候, 
将只返回大于此版本号的消息列表.


**[⬆ 回到顶部](#目录)**

4协议
======
在`protocol`目录下

* external 是对外的协议，采用`protobuf`实现
* rpc 是服务内部的调用，采用`grpc`

**[⬆ 回到顶部](#目录)**

5数据模型
======
在`doc/db`目录下


**[⬆ 回到顶部](#目录)**

6服务说明
======
进入server目录下

```shell
access
gateway
logic
manager
notify
register
```

**[⬆ 回到顶部](#目录)**

7依赖
======

### 7.1系统环境
```shell
golang >= 1.4
jdk >= 1.8 (数据处理很多服务用java编写)
```

```shell
go get -u -d github.com/golang/glog
go get -u -d github.com/coreos/etcd
go get -u -d github.com/Shopify/sarama
go get -u -d github.com/wvanbergen/kafka/consumergroup
go get -u -d github.com/tsuna/gohbase
go get -u -d github.com/garyburd/redigo/redis
go get -u -d github.com/BurntSushi/toml
go get -u -d gopkg.in/olivere/elastic.v5
go get -u -d gopkg.in/mgo.v2
go get -u -d github.com/go-sql-driver/mysql
go get -u -d github.com/satori/go.uuid
```


### 7.2第三方依赖

```shell
etcd 3.0以上版本
redis 
Mysql
kafka
HBase
ElasticSearch(可选)
```

**[⬆ 回到顶部](#目录)**

8部署
======

### 8.1普通部署

为了方便, 我们在单机上进行部署 (实际部署的时候, 每个服务角色都可以自由水平扩展)

#### 8.1.1依赖安装

* kafka安装 : http://kafka.apache.org/quickstart (默认启动即可)

安装完成之后创建两个topic:

```shell
p2p topic:
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic logic_producer_p2p

group topic:
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic logic_producer_group

```

* HBase安装 : 

安装完成之后创建:

```shell
create 'im', 'user', 'msg'
```

* redis安装 : 采用默认安装即可

* mysql安装 : 采用默认安装即可

* etcd安装 : 需要3.0以上版本, 采用默认安装即可

#### 8.1.2server安装

进入server下面的各个目录 运行 `go build`, 然后启动服务即可(因为服务做了`服务发现`, 所以对启动顺序没有要求), 这里为了简单, 每个服务我们只启动一个, 当然启动任意个都是支持的.

* gateway安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/gateway 
➜  gateway git:(master) ✗ go build
➜  gateway git:(master) ✗ ./gateway 
```

* access安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/access 
➜  access git:(master) ✗ go build
➜  access git:(master) ✗ ./access 
```

* logic安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/logic 
➜  logic git:(master) ✗ go build
➜  logic git:(master) ✗ ./logic 
```


* register安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/register 
➜  register git:(master) ✗ go build
➜  register git:(master) ✗ ./register 
```

* notify安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/notify
➜  notify git:(master) ✗ go build
➜  notify git:(master) ✗ ./notify
```

* manager安装

```shell
➜  FishChatServer2 git:(master) ✗ cd server/manager 
➜  manager git:(master) ✗ go build
➜  manager git:(master) ✗ ./manager 
```

#### 8.1.3job安装

* msg-job安装(Java)

```shell
➜  FishChatServer2 git:(master) ✗ cd jobs/msg-job
➜  msg-job git:(master) ✗ mvn clean package  -Dmaven.test.skip=true
➜  msg-job git:(master) ✗ java -jar  msg-job-core/target/msg-job-core-1.0-SNAPSHOT.jar
```

#### 8.1.4中间件服务安装

* idgen安装

```shell
➜  FishChatServer2 git:(master) ✗ cd service/idgen 
➜  idgen git:(master) ✗ go build
➜  idgen git:(master) ✗ ./idgen 
```


### 8.2容器部署

部署完全采用`Kubernetes + Docker`

所以第一步需要搭建`Kubernetes`和`Docker`, 幸运的是现在网络上已经有大量的资料了, 这块我就不多写了.

**[⬆ 回到顶部](#目录)**


9测试
======
### 9.1点对点聊天测试
进入client/p2p目录, 用户可以启动两个以上的进程, 两两之间互相聊天

```shell
➜  FishChatServer2 git:(master) ✗ cd client/p2p 
➜  p2p git:(master) ✗ go build
➜  p2p git:(master) ✗ ./p2p 
输入我的id :321
输入对方的id :收到点对点消息: 返回码[0], 对方ID[321], 消息内容[hello]
```

### 9.2群聊测试


**[⬆ 回到顶部](#目录)**

10监控
======

**[⬆ 回到顶部](#目录)**