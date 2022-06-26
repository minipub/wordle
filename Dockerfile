FROM kakie/golang:latest as BaseBuilder
COPY . /opt/wordle/
WORKDIR /opt/wordle
RUN go build

FROM ubuntu:22.04
COPY --from=BaseBuilder /opt/wordle/wordle /usr/local/bin/
CMD ["/usr/local/bin/wordle", "puzzle", "--mode", "2"]
