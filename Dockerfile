FROM golang:1.19
# Add a work directory
WORKDIR /usr/src/app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
# Copy app files
COPY . .
# Build app
RUN go build
# Expose port
EXPOSE 3000
# Start app
CMD ./star_wars