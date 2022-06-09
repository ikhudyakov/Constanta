FROM golang:1.18
WORKDIR /Constanta
COPY . .
RUN go get -u github.com/lib/pq
RUN go get -u github.com/gorilla/mux
RUN go build -o cmd/
CMD ["/main"]