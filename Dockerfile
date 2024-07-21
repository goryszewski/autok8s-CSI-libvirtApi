FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csi-libvirt ./cmd/main.go

FROM debian:latest

WORKDIR /app

RUN apt update ; apt install cryptsetup -y


COPY --from=builder /app/csi-libvirt /app/csi-libvirt

ENTRYPOINT ["/app/csi-libvirt"]