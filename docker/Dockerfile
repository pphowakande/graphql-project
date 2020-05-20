FROM registry.bookmyshow.org/golang:1.12.9 as build-base

COPY . /go/src/api/

WORKDIR /go/src/api/

RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o /api *.go

FROM registry.bookmyshow.org/alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build-base /api /home

COPY src/config/ /home/src/config

WORKDIR /home/

EXPOSE 80
EXPOSE 5001

CMD ["sh","-c","./api"]