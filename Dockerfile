FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /bin/app

FROM alpine:latest

WORKDIR /

COPY --from=build /bin/app /bin/app

RUN addgroup --gid 1000 app && \
    adduser -D -u 1000 -G app app && \
    chmod u+x /bin/app && \
    chown 1000:1000 /bin/app

EXPOSE 80 9360

ENTRYPOINT ["/bin/app"]
