FROM golang:1.13-alpine
RUN apk add --no-cache gcc musl-dev
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build
EXPOSE 8000
