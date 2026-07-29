FROM golang:alpine as base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go mod download

FROM base as build

RUN apk add build-base

RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go

FROM base as test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch as test

ARG _outputdir="/tmp/coverage"

COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

FROM alpine as final

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/docs /app/docs

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app/bookmark_service"]