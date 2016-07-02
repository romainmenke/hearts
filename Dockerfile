FROM golang:1.6.2
EXPOSE 50051
EXPOSE 8080

COPY ./ /go/src/github.com/romainmenke/hearts/

RUN go get github.com/romainmenke/hearts/...
RUN go install github.com/romainmenke/hearts

RUN git config --global push.default simple
RUN git config --global user.name heartsbot
RUN git config --global user.email romainmenke+heartsbot@gmail.com

WORKDIR /go/src/github.com/romainmenke/hearts/
CMD hearts
