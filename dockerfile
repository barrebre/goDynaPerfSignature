FROM golang:alpine
WORKDIR /go/src/barrebre/goDynaPerfSignature
ADD . /go/src/barrebre/goDynaPerfSignature
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]