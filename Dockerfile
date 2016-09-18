FROM debian:jessie

ADD gateway /gateway

WORKDIR /

RUN chmod +x /gateway

RUN ls -la / 

CMD ["/gateway"]