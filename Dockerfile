FROM index.alauda.cn/library/golang:1.8

COPY . /go/src/gateway

RUN cd /go/src/gateway \
    && go install \
    && cp -r client /go/bin/ \
    && chmod +x /go/bin/gateway /go/src/gateway/compile.sh && exit 1

WORKDIR /go/bin/

CMD ["/go/src/gateway/compile.sh"]
