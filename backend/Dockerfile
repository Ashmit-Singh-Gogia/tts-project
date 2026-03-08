# Start from a lightweight Linux image with Go installed
FROM golang:1.25-alpine

# Install the poppler-utils package (this gives us pdftotext)
RUN apk add --no-cache poppler-utils

# Set the working directory inside the container
WORKDIR /app

# Copy your Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
