FROM golang:alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN apk add build-base

RUN go mod download && \
GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go

FROM alpine

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/docs /app/docs

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app/bookmark_service"]