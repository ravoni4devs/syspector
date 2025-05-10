FROM golang:alpine
COPY . $GOPATH/src/github.com/ravoni4devs/syspector
ENV GOPROXY off
WORKDIR $GOPATH/src/github.com/ravoni4devs/syspector/cmd/example
RUN apk add --update git \
 && go build -o example get_docker_info.go \
 && mv example /usr/local/bin
CMD ["example"]
