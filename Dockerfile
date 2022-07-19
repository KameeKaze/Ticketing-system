FROM golang:alpine3.13 as build
WORKDIR /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on

COPY . .
RUN go mod download

RUN go build -o /app/main .

FROM alpine:3.16 as final

COPY --from=build /app/main /
CMD ["/main"]