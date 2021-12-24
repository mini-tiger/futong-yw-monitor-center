FROM alpine:3.13.6
RUN mkdir -p  /etc/prometheus/alertmanager && mkdir -p  /home/monitor/alertmanager
WORKDIR /home/monitor
COPY futong-yw-monitor-center .
ENTRYPOINT ["./futong-yw-monitor-center"]
