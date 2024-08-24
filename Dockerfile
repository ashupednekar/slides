
FROM golang
RUN useradd -ms /bin/sh myuser
WORKDIR /app
RUN go install golang.org/x/tools/cmd/present@latest
COPY . /app
RUN chown -R myuser:myuser /app
USER myuser
EXPOSE 3999
CMD ["present", "-http", ":3999"]
