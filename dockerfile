FROM golang:latest
RUN mkdir /symmetrical-spoon
ADD ./symmetrical-spoon /symmetrical-spoon/
WORKDIR /symmetrical-spoon
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/mattn/go-sqlite3"
RUN go get "github.com/google/uuid"
RUN go build -o /symmetrical-spoon/main .
CMD ["/symmetrical-spoon/main"]
