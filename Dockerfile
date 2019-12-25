FROM golang:1.13-alpine AS builder
RUN apk add --no-cache gcc musl-dev
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build

FROM alpine:3.11.2 AS production
COPY --from=builder /app .
EXPOSE 8000
