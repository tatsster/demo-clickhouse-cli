FROM golang:1.16.9-buster as builder

RUN apt update && apt install curl unzip -y

ENV GOPATH=/go

RUN go get -u github.com/google/wire/cmd/wire

WORKDIR /src

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download || true

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build docker
RUN make di static

######## Start a new stage from scratch #######
FROM alpine:3.13

RUN apk --no-cache add ca-certificates tzdata htop tini bash curl busybox-extras

# Change timezone to Asia/Ho_Chi_Minh
RUN rm -rf /etc/localtime\
    && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /src/out/cli /bin/
COPY --from=builder /src/out/airflow /bin/

COPY --from=builder /src/.env /

# Expose ports
EXPOSE 8081 8082

# Command to run the executable
CMD ["tini", "--"]