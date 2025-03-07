FROM debian:latest

ENV VERSION v0.3.333
RUN apt update && apt install -y wget unzip
RUN wget https://github.com/tetsu040e/zippia/releases/download/${VERSION}/zippia-${VERSION}-linux-amd64.zip
RUN unzip zippia-${VERSION}-linux-amd64.zip
RUN rm zippia-${VERSION}-linux-amd64.zip

ENTRYPOINT ["./zippia"]
CMD ["--host", "0.0.0.0"]
