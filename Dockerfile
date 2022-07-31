FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY pkg ./pkg/
COPY web ./web/
RUN go build -o pt cmd/app/main.go

# ENV PT_BACKUP_SAS="https://ptstore0.blob.core.windows.net/db-backups/?sv=2018-11-09&sr=c&st=2022-07-31&se=2023-07-31&sp=racwdl&spr=https&rscc=max-age%3D5&rscd=inline&rsce=deflate&rscl=en-US&rsct=application%2Fjson&sig=PZ9SZ1ts1ZKIA%2F3W5ftbPNWTe%2FJDpvEPccXqq8kvuKo%3D"
ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

CMD [ "./pt" ]
