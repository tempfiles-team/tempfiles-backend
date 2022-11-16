FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go build -o tempfiles-backend .

WORKDIR /dist

RUN cp /app/tempfiles-backend .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /dist/tempfiles-backend .

EXPOSE 5000

CMD ["./tempfiles-backend"]