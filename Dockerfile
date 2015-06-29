FROM golang:1.4

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

ENV APP_PORT 3030
ENV APP_LOGLEVEL info
ENV APP_FLEETENDPOINT unix:///var/run/fleet.sock

CMD ["go-wrapper", "run"]
