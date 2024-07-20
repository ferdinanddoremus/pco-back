FROM golang:1.21 AS build

WORKDIR /go/src/pco

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

# download and unzip PocketBase
FROM alpine:latest
WORKDIR /pb

COPY --from=build /go/src/pco/app .
COPY --from=build /go/src/pco/ .

EXPOSE 5000

# start PocketBase
CMD ["./app", "serve", "--http=0.0.0.0:5000"]
