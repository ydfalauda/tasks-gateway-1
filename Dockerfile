FROM index.alauda.cn/library/debian:jessie

ADD gateway /gateway

WORKDIR /

RUN chmod +x /gateway

RUN ls -la / 

CMD ["/gateway"]
