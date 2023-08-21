FROM golang:1.20-alpine3.18 AS builder

# Install make file deps
RUN apk add --no-cache git gcc musl-dev make

WORKDIR /app
COPY . .

RUN make build 

# dev image 
FROM alpine:latest

RUN apk add --no-cache gcc musl-dev make

WORKDIR /app
COPY --from=builder /app/main .
COPY ./db/migrations ./db/migrations
COPY ./scripts/wait-for /bin/wait-for

EXPOSE 8080
# command: sh -c 'echo "waiting for pg-bank:5432..." && wait-for -t 3 pg-bank:5432 -- /app/main || echo "Failed to connect to pg-banks:5432"'
CMD ["/app/main"]


# prod binary runner
# FROM scratch AS prod
# WORKDIR /app

# COPY --from=builder /app/main .
# COPY prod.env .env

# EXPOSE 8080
# CMD ["/app/main"]