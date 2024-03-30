FROM golang:1.22-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download && go mod verify

WORKDIR /app/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o music_api .

FROM alpine:3.19.1

WORKDIR /app

COPY --from=build /app/cmd/app/music_api .
COPY .env .

RUN apk add --update bash

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/81b1373f17855a4dc21156cfe1694c31d7d1792e/wait-for-it.sh /usr/bin/wait-for-it
RUN chmod +x /usr/bin/wait-for-it

EXPOSE 8083

CMD ["wait-for-it", "mariadb:3306", "--", "./music_api"]
