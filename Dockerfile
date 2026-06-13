FROM golang:1.25 AS development

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/air-verse/air@latest

RUN CGO_ENABLED=0 GOOS=linux go build -o main

CMD ["air"]

FROM alpine AS production

WORKDIR /app

COPY --from=development /app/main /app/main

CMD ["/app/main"]
