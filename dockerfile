FROM golang:alpine
WORKDIR /go/src/github.com/barrebre/goDynaPerfSignature
ADD . /go/src/github.com/barrebre/goDynaPerfSignature
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]