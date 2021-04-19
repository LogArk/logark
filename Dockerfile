FROM scratch
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY dist/logark /
EXPOSE 80
VOLUME ["/tmp"]
ENTRYPOINT ["/logark"]