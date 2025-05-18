FROM golang:1.23.5 as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o app -ldflags "-X 'main.version=1.0.0'" ./cmd/aletheia

FROM alpine:3.18

WORKDIR /root/
RUN touch checkfile
COPY --from=builder /src/app .
CMD ["./app"]
