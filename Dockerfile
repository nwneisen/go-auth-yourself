FROM golang:latest

WORKDIR /app

COPY ./ /app

# RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o go-proxy-yourself ./cmd/go-proxy-yourself/main.go" --directory=/app --recursive=true -log-prefix=false --command=/app/go-proxy-yourself
