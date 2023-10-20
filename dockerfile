# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
# WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8009

# Command to run the executable
CMD ["./main"]


# To run this application as a docker image-
# docker build -t go-auth-app . // To create a docker image for this application
# docker images     // check if the image was created successfully
# docker run -p 9000:9000 go-auth-app    // run the docker image to get a docker container and you will find your application running
# docker rmi --force <image_id1> <image_id2> // delete multiple docker images forecfully
# docker rmi <image_name>:<tag> // delete docker image using image name and tag