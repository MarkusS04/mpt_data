FROM golang:latest AS build

WORKDIR /app
COPY . .
RUN go get
RUN go build -o bin/mpt main.go

FROM debian:latest AS publish

WORKDIR /app

COPY --from=build /app/bin/mpt /app/mpt
RUN chmod +x /app/mpt

EXPOSE 8080

CMD ["/app/mpt"]