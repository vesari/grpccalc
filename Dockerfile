FROM alpine:3.11

WORKDIR /app
COPY bin/grpccalc bin/grpccalc

ENTRYPOINT ["/app/bin/grpccalc"]
