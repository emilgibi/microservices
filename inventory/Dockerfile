FROM golang:1.19-alpine

WORKDIR /app 

COPY . . 

RUN go build main.go

EXPOSE 8082

CMD ["./main"]
