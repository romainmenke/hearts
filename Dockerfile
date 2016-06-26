FROM golang:onbuild
EXPOSE 50051

RUN git config --global push.default simple
RUN git config --global user.name heartsbot
RUN git config --global user.email romainmenke+heartsbot@gmail.com
