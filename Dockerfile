FROM golang:latest
WORKDIR /todolist
COPY . .
RUN go build -o main
CMD ["./main"]
