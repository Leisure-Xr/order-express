# syntax=docker/dockerfile:1

FROM golang:1.25.4-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl && adduser -D -u 10001 app

USER app
WORKDIR /app

COPY --from=build /out/app /app/app

EXPOSE 3000

ENTRYPOINT ["/app/app"]
CMD ["serve"]
