FROM scratch
ADD ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
ADD ca-bundle.trust.crt /etc/ssl/certs/ca-bundle.trust.crt

COPY ./bin/prometheus-container-monitor /

ENTRYPOINT ["/prometheus-container-monitor"]
