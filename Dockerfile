FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o tekton-resolver .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/tekton-resolver .

USER 65532:65532

ENTRYPOINT ["./tekton-resolver"]
