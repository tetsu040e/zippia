FROM golang:latest

RUN go install github.com/tetsu040e/zippia@latest

ENTRYPOINT ["zippia"]
CMD ["--host", "0.0.0.0"]
