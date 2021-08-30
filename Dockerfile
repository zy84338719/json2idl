FROM node:12.16.1

WORKDIR /web/
COPY web .

RUN yarn && yarn build

FROM golang:alpine
WORKDIR /go/src/server
COPY server .
RUN rm -rf /go/src/server/web

COPY --from=0 /web/dist/  /go/src/server/web

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build  -ldflags "-s -w" -o server .

FROM alpine:latest
LABEL MAINTAINER="zhangyi@murphyyi.com"

WORKDIR /go/src/server

COPY --from=1 /go/src/server ./

EXPOSE 59990

ENTRYPOINT ./server