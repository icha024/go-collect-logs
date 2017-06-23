# Go Collect Logs

Lightweight tool to collect syslogs (RFC3164 & RFC5424) then display them on a web console for easy viewing/searching.

## Piping Docker logs to syslog
https://docs.docker.com/engine/admin/logging/overview/

## Build and run
To build
```
go get github.com/icha024/go-collect-logs
cd go-collect-logs
go build
```

To run, just `./go-collect-logs`
```
Usage of ./go-collect-logs:
  -host string
    	Service host to listen on. (default "0.0.0.0")
  -log-read-inteval int
    	Interval, in seconds, to read syslog into memory. (default 3)
  -max-filter int
    	Maximum number of fitlered log entries to return. (default 1000)
  -max-log int
    	Maximum number of log entries to keep. Approx 1KB/entry. (default 50000)
  -port int
    	Service port to listen on. (default 3000)
  -sev
    	Parse the syslog severity header
  -stdout
    	Print syslog received to stdout (default true)
  -syslog-host string
    	Syslog host to listen on. (default "0.0.0.0")
  -syslog-port int
    	Syslog port to listen on. (default 10514)
```

## To run in Docker:
```
docker build -t go-collect-logs
docker run --rm -ti -p 3000:3000 -p 10514:10514/UDP -p 10514:10514/TCP go-collect-logs
```

## Sending logs to go-collect-logs:
```
echo 'test message' 2>&1 | logger -P 10514 -n localhost
```

or redirect from your existing Docker container:
```
docker run --rm -ti --log-driver=syslog --log-opt syslog-format=rfc3164 --log-opt syslog-address=udp://127.0.0.1:10514 alpine /bin/sh -c "while true; do echo My Message \$RANDOM; sleep 1; done;"
```

## View logs on web console
Open `http://127.0.0.1:3000`