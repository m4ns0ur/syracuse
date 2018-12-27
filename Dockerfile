FROM alpine

RUN apk add --update ca-certificates

COPY bin/syracuse /usr/bin/syracuse

ENV POSTGRES_DSN ""
EXPOSE 8001

ENTRYPOINT syracuse -postgres-dsn=$POSTGRES_DSN