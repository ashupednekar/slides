
FROM golang:1.19.0-alpine as build-env

RUN mkdir /app
WORKDIR /app
RUN addgroup -g 10014 choreo && \
    adduser --disabled-password --no-create-home --uid 10014 --ingroup choreo choreouser
RUN go install golang.org/x/tools/cmd/present@latest

FROM scratch

COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /go/bin/present /usr/local/bin/present
COPY . /app

USER 10014

ENTRYPOINT ["/usr/local/bin/present"]
