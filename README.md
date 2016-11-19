# FishChatServer2

说明
======
吸取了第一版的经验以及我们商业版的探索。第二版我更多的思考在不要过多的实现轮子，这个版本将很多业务无关的代码去掉，全面拥抱
Kubernetes + Docker + grpc 来实现业务之上的东西。


dependence
======

<pre><code>grpc :
go get -u github.com/golang/protobuf
go get -u google.golang.org/grpc

glog:
go get -u github.com/golang/glog
</code></pre>