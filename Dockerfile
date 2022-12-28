# Step 1: Modules caching
FROM golang:1.19.4-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.19 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -o tempfiles-backend .

# GOPATH for scratch images is /
FROM scratch
WORKDIR /app
COPY --from=builder /app/tempfiles-backend .
EXPOSE 5000
CMD ["./tempfiles-backend"]