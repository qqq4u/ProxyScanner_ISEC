FROM golang:1.19 AS builder
WORKDIR /build
COPY . .
RUN go build cmd/main/main.go
FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install postgresql-12 -y
USER postgres
COPY ./CreateTable.sql .
RUN service postgresql start && \
      psql -c "ALTER USER postgres WITH PASSWORD 'password';" && \
      createdb -O postgres proxyDB && \
      psql -d proxyDB < ./CreateTable.sql && \
      service postgresql stop

USER root

WORKDIR /proxy
COPY --from=builder /build/main .

COPY . .

EXPOSE 8080
EXPOSE 8000
EXPOSE 5432

ENV PROXY_PORT=8080
ENV REPEATER_PORT=8000
ENV DB_USER=userProxy
ENV DB_NAME=proxyDB

CMD service postgresql start && ./main
