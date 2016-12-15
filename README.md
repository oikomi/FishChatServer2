# FishChatServer2

说明
======
吸取了第一版的经验以及我们商业版的探索. 第二版我更多的思考在不要过多的实现轮子, 这个版本将很多业务无关的代码去掉, 并且尽量靠拢`微服务`.
部署方式可以支持：
> * 普通部署
> * 容器部署 


架构
======

![](./doc/architecture.png)


协议
======
在`protocol`目录下

* external 是对外的协议，采用`protobuf`实现
* rpc 是服务内部的调用，采用`grpc`


服务说明
======
进入server目录下

```shell
access
gateway
logic
manager
register
```

依赖
======
```shell
etcd
redis
mongodb
kafka
ElasticSearch(可选)
```


普通部署
======



容器部署
======
部署完全采用`Kubernetes + Docker`

所以第一步需要搭建`Kubernetes`和`Docker`, 幸运的是现在网络上已经有大量的资料了, 这块我就不多写了.