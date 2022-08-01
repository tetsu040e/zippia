FROM debian:latest

RUN apt update && apt install -y wget unzip
RUN wget https://github.com/tetsu040e/zippia/releases/download/v0.3.10/zippia-v0.3.10-linux-amd64.zip
RUN unzip zippia-v0.3.10-linux-amd64.zip
RUN rm zippia-v0.3.10-linux-amd64.zip

ENTRYPOINT ["./zippia"]
CMD ["--host", "0.0.0.0"]
