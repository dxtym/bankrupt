# Build stage
FROM golang:1.22.3-alpine AS build

WORKDIR /app

COPY . .

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

RUN go build -o main main.go

# Run stage
FROM alpine

WORKDIR /app

COPY --from=build /app/main .

COPY --from=build /app/migrate ./migrate

COPY app.env .

COPY start.sh .

COPY wait-for.sh .

COPY db/migration ./migration

EXPOSE 6969

CMD [ "/app/main" ]

# Add entrypoint (before CMD to be able to override it)
ENTRYPOINT [ "/app/start.sh" ]