FROM golang:1.21

# Set the current working directory inside the container
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

# Copy files needed for installing dependencies only, so that this and the
# dependency download can be cached separately the rest of the app files.
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the workspace
COPY . .

# Build the app
# Need to use an absolute path here (/main) so that the binary is built outside
# of the working directory. Otherwise the compose volume mount will overwrite this binary.
RUN go build -o /main ./cmd/server

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
# CMD ["/main"]
CMD ["air", "-c", ".air.toml"]