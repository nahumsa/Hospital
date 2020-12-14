# Choosing a build
FROM golang:1.15
# Creating folder
RUN mkdir -p /hospital
# Defining working repository
WORKDIR /hospital
# Adding files for the container
COPY go.mod .
COPY go.sum .
COPY src/ /hospital
# Downloading dependencies
RUN go mod download
RUN go mod vendor
# Setting up enviroment variables
ENV GIN_MODE=release
ENV PORT=8080
# Building the binary
RUN go build -o hospitalbin
# Exposing the port
EXPOSE $PORT
# Running app
CMD ["./hospitalbin"]