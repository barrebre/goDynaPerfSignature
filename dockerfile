FROM golang:alpine
WORKDIR /go/src/goDynaPerfSignature
ADD . /go/src/goDynaPerfSignature
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]