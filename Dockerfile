FROM alpine

RUN apk add --update ca-certificates

COPY bin/syracuse /usr/bin/syracuse

EXPOSE 8001

ENTRYPOINT ["syracuse"]