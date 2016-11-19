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
RUN go get -u -v github.com/garyburd/redigo/redis

RUN go get -u -v github.com/oikomi/FishChatServer2

RUN cd "$GOPATH/src/github.com/oikomi/FishChatServer2/server/register" &&  go build
RUN cp "$GOPATH/src/github.com/oikomi/FishChatServer2/server/register/register" "/bin/"
RUN cp "$GOPATH/src/github.com/oikomi/FishChatServer2/server/register/register.toml" "/etc/"

EXPOSE 23000
ENTRYPOINT ["/bin/register", "-conf", "/etc/register.toml"]