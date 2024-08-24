
# Use the official Go image from the Docker Hub
FROM golang

# Set the working directory
WORKDIR /app

# Install Go Present
RUN go install golang.org/x/tools/cmd/present@latest

# Copy your presentation files into the container
# (Assuming your Go Present files are in the same directory as this Dockerfile)
COPY . /app

# Expose port 3999 for the presentation server
EXPOSE 3999

# Command to run Go Present
CMD ["present", "-http", ":3999"]
