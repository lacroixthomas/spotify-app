FROM golang:1.16.2
RUN mkdir /playlist
WORKDIR /playlist
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go test -v
RUN go build main.go
EXPOSE 8080
ENTRYPOINT [ "/playlist/main" ]