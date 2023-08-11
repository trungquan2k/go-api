FROM golang:1.20.5

WORKDIR /app

ENV config=docker

COPY .. /app

RUN go mod download

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

EXPOSE 9091

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command=./main
