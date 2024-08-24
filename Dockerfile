
FROM golang:1.19.0-alpine as build-env

RUN mkdir /app
WORKDIR /app
RUN addgroup -g 10014 choreo && \
    adduser --disabled-password --no-create-home --uid 10014 --ingroup choreo choreouser
RUN go install golang.org/x/tools/cmd/present@latest

COPY . /app

USER 10014

ENTRYPOINT ["present"]
