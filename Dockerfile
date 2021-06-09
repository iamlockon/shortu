FROM golang:1.15-alpine AS build-stage

ARG PORT=8080
ARG MAINDIR=./cmd/

WORKDIR /app
COPY . .
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main $MAINDIR

FROM alpine:3.12
COPY --from=build-stage /app/main /app/main
EXPOSE $PORT
RUN apk --no-cache add curl
CMD ["/app/main"]
