FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build cmd/online-buddy/main.go


FROM golang:1.22-alpine

ARG API_PORT
ENV API_PORT $API_PORT

ARG REDIS_CONNECTION_URL
ENV REDIS_CONNECTION_URL $REDIS_CONNECTION_URL

WORKDIR /app

# https://stackoverflow.com/a/50861580
RUN apk add --no-cache libc6-compat 

COPY --from=builder /app/main /app/main
COPY web web

EXPOSE $API_PORT

CMD ["/app/main"]