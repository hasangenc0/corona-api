FROM golang:1.13.1 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./

FROM alpine:latest AS production

COPY --from=builder /app .
ENV ENV="development"

CMD ["./main"]