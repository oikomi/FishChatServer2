FROM golang

MAINTAINER miaohong@miaohong.org

ENV GOPATH /go
ENV USER root

RUN mkdir -p "$GOPATH/src/" "$GOPATH/bin"
RUN go get -u -v github.com/BurntSushi/toml
RUN go get -u -v github.com/golang/glog
RUN go get -u -v github.com/golang/protobuf/proto
RUN go get -u -v github.com/coreos/etcd/clientv3
RUN go get -u -v gopkg.in/mgo.v2
RUN go get -u -v github.com/Shopify/sarama
RUN go get -u -v github.com/wvanbergen/kafka/consumergroup
RUN go get -u -v gopkg.in/olivere/elastic.v5

RUN go get -u -v github.com/oikomi/FishChatServer2

RUN cd "$GOPATH/src/github.com/oikomi/FishChatServer2/server/logic" &&  go build
RUN cp "$GOPATH/src/github.com/oikomi/FishChatServer2/server/logic/logic" "/bin/"
RUN cp "$GOPATH/src/github.com/oikomi/FishChatServer2/server/logic/logic.toml" "/etc/"

EXPOSE 21000
ENTRYPOINT ["/bin/logic", "-conf", "/etc/logic.toml"]