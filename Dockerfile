FROM golang:1.7

RUN mkdir -p /go/src/github.com/elcct/taillachat

COPY . /go/src/github.com/elcct/taillachat

RUN go get -u github.com/elcct/taillachat

EXPOSE 8000

CMD taillachat
