FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /bin/app

FROM alpine:latest

WORKDIR /

COPY --from=build /bin/app /bin/imaginarium
COPY ./config.yml /etc/imaginarium/config.yml

RUN addgroup --gid 1000 imaginarium && \
    adduser -D -u 1000 -G imaginarium imaginarium && \
    chmod u+x /bin/imaginarium && \
    chown imaginarium:imaginarium /bin/imaginarium && \
    chown imaginarium:imaginarium -R /etc/imaginarium && \
    chmod 0660 /etc/imaginarium/config.yml

USER imaginarium
EXPOSE 80 9360

ENTRYPOINT ["/bin/imaginarium", "-config", "/etc/imaginarium/config.yml"]
