FROM golang:1.4.2-onbuild

ADD Dockefile.run /go/bin/Dockerfile

WORKDIR /go/bin

CMD tar -cf - .
