FROM golang:1.25.0-alpine

WORKDIR /workdir

COPY ./ /workdir/
COPY ./views /apps/views/
COPY ./static /apps/static/

RUN go build -o /apps/SimpleHomeInventory

WORKDIR /apps/

RUN rm -rf /workdir

EXPOSE 7070/tcp

ENTRYPOINT [ "/apps/SimpleHomeInventory" ]