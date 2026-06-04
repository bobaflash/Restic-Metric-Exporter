FROM golang:1.26.4-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /data

RUN echo Building for linux
RUN mkdir -p bin
COPY . .
RUN go get -v -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -tags timetzdata -o bin/app -a ./cmd/main.go

FROM scratch

WORKDIR /data
COPY --from=builder /data/bin/app /data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ="Europe/Berlin"
CMD [ "/data/app" ]