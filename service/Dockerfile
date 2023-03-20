FROM golang:1.19

EXPOSE 80

COPY . /app
WORKDIR /app

RUN ["go", "get", "github.com/lib/pq"]
RUN ["go", "build", "-o", "executable", "cmd/main.go"]

ENTRYPOINT ["/app/executable"]
