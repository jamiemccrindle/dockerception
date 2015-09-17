# Dockerception

## Docker building dockers - keeping them small

### 1.8 Permissions error

Docker introduced an issue in 1.8 that affects how dockerception works
[https://github.com/docker/docker/issues/15785](https://github.com/docker/docker/issues/15785):

    unable to prepare context: unable to extract stdin to temporary context directory: lchown ...: operation not permitted

I'm using the following as a workaround:

    > cat ~/bin/dockerception
    set -e
    BUILD_DIR=`mktemp -d /tmp/dockerception-$1.XXXXXX`
    echo $BUILD_DIR
    docker build -t $1-builder .
    docker run $1-builder > $BUILD_DIR/$1.tar
    tar -C $BUILD_DIR -xvf $BUILD_DIR/$1.tar
    docker build -t $1 $BUILD_DIR
    rm -r $BUILD_DIR

### tl;dr

You can split out your docker build process into a 'builder' docker and a 'runtime' docker to keep your docker runtime images
as small as possible. This repository is an example of that. To build the runtime docker image, clone this project and then
run the following command:

    docker build -t builder .; docker run builder | docker build -t dockerception -

### The longer version

Having an entirely self contained build process within docker is convenient. A downside is that doing this often means
that there are build time dependencies that are carried over to your runtime e.g. the official golang builder docker
weighs in at 514.8mb before you even add your project in. A better solution would be to be able to build a 'builder'
docker image and then use that to construct a 'runtime' docker image.

There is a proposal for nested builds: Proposal: Nested builds #7115, but it has been open for some time.

I recall seeing the following feature in the release notes for Docker 1.1.0:

> Allow a tar file as context for docker build
>
> You can now pass a tar archive to docker build as context. This can be used to automate docker builds, for example: cat context.tar | docker build - or docker run builder_image | docker build -

but I hadn't seen an examples of it being used, so I decided to try it out.

Skipping to the end, here is the line that builds our builder docker image and then builds the final runtime docker image:

    docker build -t builder .; docker run builder | docker build -t dockerception -

which does the following:

* Builds a 'builder' docker image using the Dockerfile in the current directory (docker build -t builder .)
* Runs the 'builder' docker which builds the sources in the current directory and outputs them as a tar stream (docker run builder)
* Builds an image called 'dockerception' from the tar stream which contains a Dockerfile and the binary (docker build -t dockerception -)

The Dockerfile for the builder looks as follows:

    FROM golang:1.4.2-onbuild

    # Add the runtime dockerfile into the context as Dockerfile
    ADD Dockerfile.run /go/bin/Dockerfile

    # Set the workdir to be /go/bin which is where the binaries are built
    WORKDIR /go/bin

    # Export the WORKDIR as a tar stream
    CMD tar -cf - .

The Dockerfile for the runtime image looks as follows:

    FROM flynn/busybox

    # Add the binary
    ADD app /bin/app

    EXPOSE 8001

    # Run the /bin/app by default
    CMD ["/bin/app"]

Resulting in a 10.53 MB docker image. It should be possible to build the runtime docker image using scratch instead of
busybox but I'll leave that as an exercise for the reader.

If this still seems confusing, here's a deeper dive into what's happening:

'docker build' usually points at a directory e.g. 'docker build .'. This directory is known as the docker build context.
The directory typically has a Dockerfile in it and any other resources you want to add via the COPY command e.g. for a
golang project your directory / context may look something like this:

    > tree .
    .
    |-- Dockerfile
    |-- main.go

And your docker build command would look as follows:

    docker build .

Now, it turns out that docker can also accept that context via a tar file as in, if you run:

    > tar -cvf /tmp/image.tar .
    a .
    a ./Dockerfile
    a ./main.go

You could run docker build as follows:

    > cat /tmp/image.tar | docker build . -

All that we need to do now is create a 'builder' docker that can construct a context directory and write it to standard out as a tar stream that we can pipe to docker build

This is also available on my blog [jamie.mccrindle.org](http://jamie.mccrindle.org/2015/04/dockerception-how-to-have-docker-build.html)