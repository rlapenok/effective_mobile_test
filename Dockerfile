# Stage 1: Build stage
FROM golang:1.22-alpine AS build

# Set the working directory
WORKDIR /app

COPY . .
RUN go mod download
RUN cd cmd && go build -o ../effective_mobile_test main.go
CMD [ "./effective_mobile_test","--path_to_config", ".env" ]