# Go Collect Logs

## Piping Docker logs to syslog
https://docs.docker.com/engine/admin/logging/overview/

For example:
```
docker run --rm -ti --log-driver=syslog --log-opt syslog-format=rfc3164 --log-opt syslog-address=udp://127.0.0.1:10514 alpine /bin/sh -c "while true; do echo My Message \$RANDOM; sleep 1; done;"
```