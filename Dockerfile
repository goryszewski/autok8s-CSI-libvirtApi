FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csi-libvirt ./main.go

FROM debian:latest

WORKDIR /app

# RUN apk add --no-cache e2fsprogs

COPY --from=builder /app/csi-libvirt /app/csi-libvirt

ENTRYPOINT ["/app/csi-libvirt"]