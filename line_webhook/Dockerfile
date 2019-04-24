
FROM golang:1.12.4 as build_base

LABEL maintainer="Sarakorn <sarakorn@gmail.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)
RUN go mod download



FROM build_base as server_builder
# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/line-webhook .


# start from a fresh Alpine image, to reduce the image size and not shipp the Go compiler in production artifacts.
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=server_builder /go/bin/line-webhook .
# COPY .env .

EXPOSE 8080

CMD ["./line-webhook"]