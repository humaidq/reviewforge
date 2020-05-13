FROM golang:1.13

RUN apt update && apt upgrade -y

WORKDIR /root/reviewforge
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080:8080

CMD ["reviewforge", "run"]
