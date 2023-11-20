FROM golang:1.21
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY chat-server .
RUN go build -o main .
CMD ["./main"]
