FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o auth-app cmd/main.go

FROM alpine AS runner

RUN apk --no-cache add curl

WORKDIR /app

COPY --from=build /app/auth-app ./
COPY --from=build /app/config ./config/

CMD ["./auth-app"]