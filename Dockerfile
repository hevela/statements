# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/hevela/statements

# our working dir will be the root of the service
WORKDIR /go/src/github.com/hevela/statements

# enable gooneeleven flag
ENV GO111MODULE=on

# Download and install dependencies
RUN go mod download
RUN go mod vendor

# build the service
RUN go build -o ./. cmd/app/main.go

ENTRYPOINT ["./main"]