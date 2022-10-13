FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/auth-service
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/auth-service /go/src/auth-service


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/auth-service /usr/local/bin/auth-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["auth-service"]
