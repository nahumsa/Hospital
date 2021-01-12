# Choosing a build
FROM golang:1.15
# Creating folder
RUN mkdir -p /go/src/github.com/nahumsa/hospital-management/
# Defining working repository
WORKDIR /go/src/github.com/nahumsa/hospital-management/
RUN mkdir -p /src/
# Adding files for the container
COPY go.mod .
COPY go.sum .
COPY src/ src/
# Downloading dependencies
RUN go mod download
RUN go mod vendor
# Setting up enviroment variables
ENV GIN_MODE=release
ENV PORT=8080
# Building the binary
WORKDIR /go/src/github.com/nahumsa/hospital-management/src/
RUN go build -o hospitalbin
# Exposing the port
EXPOSE $PORT
# Running app
CMD ["./hospitalbin"]