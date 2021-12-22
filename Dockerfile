# Build stage

FROM golang:1.16 as build

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o cmd -v

# Run stage

FROM alpine:latest

WORKDIR /app

COPY --from=build /build/cmd .

CMD ["/app/cmd"]