FROM golang:alpine as GO_BUILD
WORKDIR /go/src/github.com/88250/lute-http/
ADD . /go/src/github.com/88250/lute-http/
ENV GO111MODULE=on
RUN apk add --no-cache gcc musl-dev git && go build -i -v

FROM alpine:latest
LABEL maintainer="Liang Ding<d@b3log.org>"
WORKDIR /opt/lute-http/
COPY --from=GO_BUILD /go/src/github.com/88250/lute-http/lute-http /opt/lute-http/lute-http
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
EXPOSE 8249

ENTRYPOINT [ "/opt/lute-http/lute-http" ]
