FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o gateway ./cmd/gateway/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/gateway .

#COPY --from=builder /app/docs ./docs

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata

ARG PORT=8080
ENV PORT=$PORT

EXPOSE ${PORT}
