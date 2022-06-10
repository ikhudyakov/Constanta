FROM golang:1.18
WORKDIR /home/ilei/Downloads/Constanta/
COPY . .
COPY /cmd .
RUN go get -u github.com/lib/pq
RUN go get -u github.com/gorilla/mux
RUN go build -o main .
CMD ["./main"]
