FROM golang:1.19-alpine

WORKDIR /app

COPY . .
RUN go mod download

COPY *.go ./

RUN go build -o ./tempfile-backend

EXPOSE 5000

CMD [ "./tempfile-backend" ]