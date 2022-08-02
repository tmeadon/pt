FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY pkg ./pkg/
COPY web ./web/
RUN go build -o pt cmd/app/main.go

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

CMD [ "./pt" ]
