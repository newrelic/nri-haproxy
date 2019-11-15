FROM golang:1.10 as builder
COPY . /go/src/github.com/newrelic/nri-haproxy/
RUN cd /go/src/github.com/newrelic/nri-haproxy && \
    make && \
    strip ./bin/nri-haproxy

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-haproxy/bin/nri-haproxy /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-haproxy
COPY --from=builder /go/src/github.com/newrelic/nri-haproxy/haproxy-definition.yml /nri-sidecar/newrelic-infra/newrelic-integrations/definition.yml
USER 1000
