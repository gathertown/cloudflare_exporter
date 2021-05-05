FROM golang:1.15.7-alpine3.12 as builder
RUN apk update && apk add ca-certificates curl git make tzdata
RUN adduser -u 5003 --gecos '' --disabled-password --no-create-home gather
COPY . /app
WORKDIR /app
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/cloudflare_exporter /bin/cloudflare_exporter
COPY --from=builder /etc/passwd /etc/passwd
USER gather
CMD ["cloudflare_exporter"]
