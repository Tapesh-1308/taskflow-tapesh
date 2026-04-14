FROM golang:1.25-alpine AS builder
WORKDIR /app

# install git (needed for private deps sometimes)
RUN apk add --no-cache git

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build (FIX: point to correct main package)
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o main ./cmd

# -------- minimal runtime --------
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./main"]