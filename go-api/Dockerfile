# Specifies a parent image
FROM golang:1.22.1-alpine3.19

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go mod download

# Install CompileDaemon for hot-reloading in development
RUN go install github.com/githubnemo/CompileDaemon@latest

# Builds your app with optional configuration
RUN go build -o /app

# Tells Docker which network port your container listens on
EXPOSE 8081

# Specifies the executable command that runs when the container starts
CMD [ “/app ]
 