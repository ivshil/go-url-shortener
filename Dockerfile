FROM golang:1.20.5

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE $GOAPP_PORT

ENV PGDB_HOST=$PGDB_HOST
ENV PGDB_PORT=$PGDB_PORT
ENV PGDB_USER=$PGDB_USER
ENV PGDB_PASS=$PGDB_PASS
ENV PGDB_NAME=$PGDB_NAME

CMD ["./main"]