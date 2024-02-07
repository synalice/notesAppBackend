FROM golang:1.21

WORKDIR /usr/src/notesAppBackend

COPY go.mod go.sum  ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/notesAppBackend notesAppBackend/cmd/notesAppBackend

EXPOSE 8080

CMD ["notesAppBackend"]
