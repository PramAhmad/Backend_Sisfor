FROM golang:1.18.1-alpine
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -buildvcs=false -o main .
CMD ["/app/main"]