FROM golang:1.12 as builder

# Set Environment Variables
ENV HOME /app
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build app
RUN go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
    

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD [ "./main" ]