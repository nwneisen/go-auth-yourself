FROM golang:1.20.3-alpine3.17

# Install CompileDaemon
RUN apk add --no-cache git && \
  git clone https://github.com/githubnemo/CompileDaemon.git && \
  cd CompileDaemon && \
  go build -o CompileDaemon . && \
  chmod +x CompileDaemon && \
  mv CompileDaemon /usr/local/bin

WORKDIR /app
COPY ./ /app

ENTRYPOINT CompileDaemon --build="go build -o go-proxy-yourself ./cmd/go-proxy-yourself/main.go" --directory=/app --recursive=true -log-prefix=false --command=/app/go-proxy-yourself
