FROM golang:1.15-alpine AS build_base

RUN apk add --no-cache git
WORKDIR /tmp/filehost

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./out/filehost .

FROM alpine:3.9
WORKDIR /app

COPY --from=build_base /tmp/filehost/out/filehost /app/filehost
COPY --from=build_base /tmp/filehost/static/ /app/static/
EXPOSE 80

CMD ["/app/filehost"]
