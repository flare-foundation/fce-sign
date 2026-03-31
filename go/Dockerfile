ARG TEE_NODE_VERSION=main

FROM golang:1.25.1-alpine AS node-builder

RUN apk add --no-cache git ca-certificates tzdata

ARG TEE_NODE_VERSION
WORKDIR /node
RUN git clone --depth 1 --branch ${TEE_NODE_VERSION} https://github.com/flare-foundation/tee-node.git .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /node/server cmd/extension/main.go

FROM golang:1.25.1-alpine AS ext-builder

RUN apk add --no-cache git
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/extension-tee .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache gosu --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN adduser -DH -u 1000 server && adduser -DH -u 1001 extension

COPY --from=node-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=node-builder /node/assets/google_confidential_space_root.crt assets/google_confidential_space_root.crt
COPY --from=node-builder /node/server ./
COPY --from=ext-builder /app/extension-tee ./

ENV TZ=UTC
ENV MODE=1
ENV CONFIG_PORT=5500
ENV SIGN_PORT=8882
ENV EXTENSION_PORT=8883

LABEL "tee.launch_policy.allow_env_override"="LOG_LEVEL"

EXPOSE 5500

CMD ["sh", "-c", "./server & gosu extension ./extension-tee"]
