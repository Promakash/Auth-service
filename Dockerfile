FROM golang:1.22.3 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o auth-app cmd/main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

COPY --from=build /app/auth-app ./
COPY --from=build /app/config ./config/

CMD ["./auth-app"]