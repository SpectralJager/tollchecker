# Build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go build -o receive cmd/main.go

# Running stage
FROM alpine
WORKDIR /build
COPY --from=build /build/receive /build/receive
CMD [ "./receive" ]