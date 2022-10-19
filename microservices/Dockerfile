FROM golang:1.19

WORKDIR /app

COPY ./ /app/

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download

EXPOSE 8080

CMD make run 
