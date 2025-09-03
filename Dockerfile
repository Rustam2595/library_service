FROM golang:1.24.0

WORKDIR /app

COPY . .

RUN go build cmd/books/main.go

CMD ["./main"]  
#FROM postgres:14
#ENV POSTGRES_DB=myapp
#ENV POSTGRES_USER=myuser
##ENV POSTGRES_PASSWORD=mypassword
##COPY ./init-scripts/ /docker-entrypoint-initdb.d/
##COPY ./postgresql.conf /etc/postgresql/postgresql.conf
#EXPOSE 5432
