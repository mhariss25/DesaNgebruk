FROM golang:1.21.5

RUN apt-get update && apt-get upgrade 

RUN apt-get install curl git -y

WORKDIR /app

COPY . /app

RUN go mod tidy

RUN go compile main.go -o /app

ENTRYPOINT [ "/app" ]