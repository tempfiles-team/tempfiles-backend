FROM ubuntu:latest

RUN mkdir /app

WORKDIR /app

COPY ./tempfiles-backend ./tempfiles-backend

EXPOSE 5000

CMD ["./tempfiles-backend"]