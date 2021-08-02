# Copyright (c) 2021. Quirino Gervacio
# MIT License. All Rights Reserved

# ----------------
# 1st Stage: Build
# ----------------

FROM golang:1.16.6-alpine3.14 as builder
LABEL maintainer="Quirino Gervacio <qgervacio@gmail.com>"

ENV GOPRIVATE github.com/spolarium

WORKDIR /app

COPY . ./
RUN go mod download -x
RUN CGO_ENABLED=0 GOOS=linux go build main.go

# --------------
# 2nd Stage: Run
# --------------

FROM alpine:3.14.0
LABEL maintainer="Quirino Gervacio <qgervacio@gmail.com>"

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["/app/main"]
