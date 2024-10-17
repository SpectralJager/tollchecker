# Build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go build -o obu cmd/main.go

# Running stage
FROM alpine
WORKDIR /build
COPY --from=build /build/obu /build/obu
CMD [ "./obu" ]