FROM golang
MAINTAINER Hugo

ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /go/src/rem_cli
COPY . .
RUN CGO_ENABLED=0 go build -o rem
#EXPOSE 8080
#CMD ["./rem"]

#FROM scratch
#COPY . /
#CMD ["/rem"]
#ENTRYPOINT ["/rem"]