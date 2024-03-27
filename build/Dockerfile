FROM golang:1.22-alpine AS build

ENV CGO_ENABLED=1

WORKDIR /app
RUN apk add --update \
     gcc \
     musl-dev

COPY . .
RUN go mod download && go mod verify

WORKDIR /app/cmd/app
RUN GOOS=linux go build -o music_api .

WORKDIR /app
FROM busybox

WORKDIR /app

COPY --from=build /app/cmd/app/music_api .

EXPOSE 8083

CMD ["./music_api"]
