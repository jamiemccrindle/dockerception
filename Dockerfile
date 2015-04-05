FROM golang:1.4.2-onbuild

# Add the runtime dockerfile into the context as Dockerfile
ADD Dockerfile.run /go/bin/Dockerfile

# Set the workdir to be /go/bin which is where the binaries are built
WORKDIR /go/bin

# Export the WORKDIR as a tar stream
CMD tar -cf - .
