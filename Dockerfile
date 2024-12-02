# Step 1: Use the official Go image as the build environment
FROM golang:1.23 AS build

# Step 2: Install dependencies for `air`
RUN go install github.com/cosmtrek/air@latest

# Step 3: Set the working directory inside the container
WORKDIR /app

# Step 4: Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 5: Copy the entire source code into the container
COPY . .

# Step 6: Set the environment variable for `air` to be available in $PATH
ENV PATH=$PATH:/root/go/bin

# Step 7: Use a smaller Alpine-based image for the final container
FROM alpine:latest

# Step 8: Install necessary certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Step 9: Set the working directory for the final image
WORKDIR /root/

# Step 10: Copy the application files from the build environment
COPY --from=build /app /app

# Step 11: Expose the application's port
EXPOSE 8080

# Step 12: Set the environment variable for `air` to be available in $PATH
ENV PATH=$PATH:/root/go/bin

# Step 13: Command to run `air` for live-reloading
CMD ["air"]
