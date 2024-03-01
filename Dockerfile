ARG GO_VER
ARG ALPINE_VER

FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder

ARG APP_NAME

WORKDIR /project

COPY . .

RUN go build -o ./server ./cmd/${APP_NAME}/main.go

FROM alpine:${ALPINE_VER}

WORKDIR /project

RUN apk update && \
    apk add ca-certificates

COPY --from=builder /project .

ENTRYPOINT [ "./server" ]
