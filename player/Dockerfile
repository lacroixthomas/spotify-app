FROM golang:1.16.2
RUN mkdir /player
WORKDIR /player
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build main.go
EXPOSE 8080
ENTRYPOINT [ "/player/main" ]