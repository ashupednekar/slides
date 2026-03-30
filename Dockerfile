# present@latest tracks x/tools, which currently needs Go >= 1.25
FROM golang:1.25-alpine AS build
ENV CGO_ENABLED=0
RUN go install golang.org/x/tools/cmd/present@latest

FROM scratch
COPY --from=build /go/bin/present /present
WORKDIR /app
COPY . .
EXPOSE 3999
ENTRYPOINT ["/present", "-http", ":3999"]
