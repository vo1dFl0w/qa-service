FROM golang:1.24.3-alpine AS builder

WORKDIR /qa-service

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=configs/config.yaml
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o qa-service ./cmd/qa-service/main.go

FROM alpine AS runner
RUN apk add --no-cache ca-certificates

WORKDIR /qa-service
COPY --from=builder /qa-service/qa-service /qa-service/qa-service
COPY --from=builder /qa-service/configs/config.yaml /qa-service/configs/config.yaml

EXPOSE 8080

CMD ["./qa-service"]