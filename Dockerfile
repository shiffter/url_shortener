FROM golang:latest as builder

WORKDIR /app
ENV CGO_ENABLED=0
COPY . .

RUN go build -o . cmd/main.go

FROM alpine:latest

COPY --from=builder /app/main /app/main
COPY --from=builder /app/config /app/config
WORKDIR /app

EXPOSE 8081

RUN apk --no-cache add ca-certificates

ENTRYPOINT ["./main"]