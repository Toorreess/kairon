FROM golang:alpine AS build
WORKDIR /go/src/kairon
COPY . .
RUN export CGO_ENABLED=0 && go build -o /go/bin/kairon ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=build /go/bin/kairon /go/bin/kairon
WORKDIR /go/src/kairon
COPY config/config.yml config/config.yml
ENV GOPATH=/go
ENTRYPOINT ["/go/bin/kairon"]