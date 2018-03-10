FROM golang:1.9-alpine as server

# GOPATH = /go in the golang image
# also $GOPATH/bin has been added to path

WORKDIR /go/src/

# copy API server src to be compiled
COPY API/ ./

# currently, we need to install several dependencies
# note that we must install Docker, but when we run the container, we must
# mount the /var/run/docker.sock of the host onto the container so the host
# docker daemon spins up sibling containers (as opposed to Docker in Docker)
#
# TODO (cw|3.9.18) once we remove the dependency on gb, we will only need to
# install docker here.
RUN apk update && \
    apk add --no-cache git && \
    go get github.com/constabulary/gb/... && \
    apk del git && \
    apk add docker

# compile and install server binary within container
RUN gb build && \
    mv ./bin/compilebox /go/bin

FROM alpine

WORKDIR /bin/

# copy over single binary from build stage --^
COPY --from=server /go/bin/compilebox .

EXPOSE 6666

# run comilebox API server
CMD ["compilebox"]
