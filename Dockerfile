FROM golang:1.25.0-alpine

WORKDIR /workdir

# Make docker builds faster for go projects
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# "Normal" Docker build process
COPY . .
COPY ./views /apps/views/
COPY ./static /apps/static/

RUN go build -v -cover -o /apps/SimpleHomeInventory

WORKDIR /apps/

RUN rm -rf /workdir

EXPOSE 7070/tcp

ENTRYPOINT [ "/apps/SimpleHomeInventory" ]