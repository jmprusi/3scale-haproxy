#### Builder
FROM golang:1.12 as builder


WORKDIR /tmp/3scale-haproxy

COPY . .

RUN go build threescale_haproxy

#### Runtime
FROM registry.access.redhat.com/ubi8/ubi-minimal

ENV HOSTNAME ACCESS_TOKEN 3SCALE_ADMIN_URL SERVICE_ID

WORKDIR /root/

COPY --from=builder /tmp/3scale-haproxy/threescale_haproxy .

EXPOSE 12345

ENTRYPOINT ["./threescale_haproxy"] 
