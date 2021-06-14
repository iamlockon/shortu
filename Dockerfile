FROM golang:1.15-alpine AS build-stage

ARG PORT=8080
ARG MAINDIR=./cmd/

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main $MAINDIR
RUN go get github.com/jackc/tern@v1.12.4
RUN echo 'tern version:'
RUN tern version

FROM alpine:3.12
COPY --from=build-stage /app/main /app/scripts/wait_for.sh /go/bin/tern /app/tern.conf /app/
COPY --from=build-stage /app/sqls /app/sqls
RUN ls -l  /app
EXPOSE $PORT
RUN apk --no-cache add curl bash
CMD ["/app/main"]
