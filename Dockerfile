FROM golang:1.4.2-onbuild

ADD Dockerfile.run /go/bin/Dockerfile

WORKDIR /go/bin

CMD tar -cf - .
