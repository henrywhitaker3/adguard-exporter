FROM scratch
WORKDIR /

COPY adguard-exporter /adguard-exporter
USER 65532:65532

ENTRYPOINT ["/adguard-exporter"]
