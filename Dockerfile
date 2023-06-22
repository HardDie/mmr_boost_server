##
# Build
##
FROM golang:alpine AS builder
LABEL stage=builder
WORKDIR /app
# Install grpc files
RUN apk add --no-cache protoc make git
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# Install dependency
COPY go.mod go.sum .
RUN go mod download
# Build project
COPY . .
RUN make proto
RUN go build -o server cmd/mmr_boost_server/main.go

##
# Run image
##
FROM alpine
WORKDIR /app
COPY --from=builder /app/server ./server
COPY --from=builder /app/api.swagger.yaml ./api.swagger.yaml
CMD ./server
