FROM golang:1.9.0-alpine3.6 AS BUILD

MAINTAINER CMGS <ilskdw@gmail.com>

# make binary
RUN apk add --no-cache git curl make \
    && curl https://glide.sh/get | sh \
    && go get -d github.com/projecteru2/lambda
WORKDIR /go/src/github.com/projecteru2/lambda
RUN make build && ./lambda --version

FROM alpine:3.6

MAINTAINER CMGS <ilskdw@gmail.com>

RUN mkdir /etc/eru/
COPY --from=BUILD /go/src/github.com/projecteru2/lambda/lambda /usr/bin/lambda
COPY --from=BUILD /go/src/github.com/projecteru2/lambda/lambda.yaml.example /etc/eru/lambda.yaml.example
