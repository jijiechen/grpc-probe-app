# build command : docker build . -t jijiechen/grpc-probe-app
# run command : docker run -it jijiechen/grpc-probe-app

FROM golang:1.21-bookworm as builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o grpc-probe-app


FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/grpc-probe-app /app/grpc-probe-app
CMD ["/app/grpc-probe-app"]

