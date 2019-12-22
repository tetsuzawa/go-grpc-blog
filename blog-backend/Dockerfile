FROM golang:1.13

WORKDIR /go/src/blog-backend
COPY . .
RUN go get -u github.com/pilu/fresh
CMD ["fresh"]
EXPOSE 31090