# Build stage
FROM golang:1.22.3-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o main main.go

# Run stage
FROM alpine

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 6969

CMD [ "/app/main" ]