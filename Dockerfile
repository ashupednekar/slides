# present@latest tracks x/tools, which currently needs Go >= 1.25
FROM golang:1.25-alpine AS build
ENV CGO_ENABLED=0
RUN go install golang.org/x/tools/cmd/present@latest

FROM scratch
COPY --from=build /go/bin/present /present
WORKDIR /app
COPY . .
EXPOSE 3999
# Custom templates in present-assets/ (slides.tmpl has no closing "Thank you" slide — see present-assets/templates/slides.tmpl)
ENTRYPOINT ["/present", "-http", ":3999", "-base=/app/present-assets"]
