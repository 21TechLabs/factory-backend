# Use the official Golang image as the base
FROM golang:1.25.3-alpine
# Set the working directory to /app
WORKDIR /app
# Copy the Go source code to the working directory
COPY . .
# Install dependencies
RUN go mod tidy
# Build the Go application
RUN go build -o /app/main . && chmod +x main
# Expose the port your application will listen on (replace 8080 with your desired port)
EXPOSE 6969

# Define the command to run when the container starts
CMD ["/app/main"]