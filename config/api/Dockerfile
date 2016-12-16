# Base this docker container off the official golang docker image.
# Docker containers inherit everything from their base.
FROM golang:1.6.3

# Create a directory inside the container to store all our application and then make it the working directory.
RUN mkdir -p /go/src/github.com/pick-up-api
WORKDIR /go/src/github.com/pick-up-api

# Copy the api directory into the container.
COPY ./ /go/src/github.com/pick-up-api

# Download and install test dependencies
RUN go get github.com/stretchr/testify/assert

# Download and install any required third party dependencies into the container.
RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

# Set environment variables
ENV PORT 8080
ENV DB_USER "raccoon"
ENV DB_PASS "pickEmUp"
ENV DB_TABLE "pickup"
ENV DB_HOST "pickup-postgres"

# Expose port 3000 to the host so we can access the gin proxy
EXPOSE 3000

# Now tell Docker what command to run when the container starts
CMD gin run
