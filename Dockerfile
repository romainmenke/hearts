FROM golang:onbuild
EXPOSE 50051
EXPOSE 8080

RUN git config --global push.default simple
RUN git config --global user.name heartsbot
RUN git config --global user.email romainmenke+heartsbot@gmail.com
