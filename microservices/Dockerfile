FROM golang:1.19

WORKDIR /app

COPY ./ /app/
# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client
# make wait for postgres
RUN chmod +x wait-for-postgres.sh
# build go app
RUN go mod download

EXPOSE 8080

CMD make run 
