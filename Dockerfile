FROM ubuntu:18.04

COPY bin /opt/qeuekeeper

EXPOSE 8086

WORKDIR /opt/qeuekeeper
CMD ["./queuekeeper"]
