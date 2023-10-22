FROM golang:latest

WORKDIR /app

COPY . .

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client



# install net-streaming
RUN go install github.com/nats-io/nats-server/v2@latest

# make wait-for-postgres.sh and wait-for-nats.sh executable
RUN chmod +x wait-for-postgres.sh && chmod +x wait-for-nats.sh

# build go app
RUN go mod download
RUN go build cmd/wbapp/main.go

CMD [ "./main" ]