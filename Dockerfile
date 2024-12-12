FROM golang:latest AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY main.go /
RUN go build -o /cacheServer /main.go


FROM scratch
COPY --from=builder /cacheServer /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/cacheServer"]

