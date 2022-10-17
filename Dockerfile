FROM golang:1.18

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o paymentservice ./cmd

CMD ["./paymentservice"]