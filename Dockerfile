FROM golang:1.8

COPY . /go/src/gateway

RUN cd /go/src/gateway && go install && cp -r client /go/bin/

WORKDIR /go/bin/

RUN chmod +x /go/bin/gateway /go/src/gateway/compile.sh

RUN ls -la /go/bin/

CMD ["/go/src/gateway/compile.sh"]
