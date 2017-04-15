package main

import (
	"flag"
	"fmt"
	"github.com/icha024/go-collect-logs/sse"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"log"
	"net/http"
	// "os"
	"time"
)

func main() {
	var maxLogEntries = flag.Int("max-log-entries", 10, "Maximum number of log entries to keep. Approx 1KB/entry.")
	var logReadInteval = flag.Int("log-read-inteval", 3, "Interval, in seconds, to read syslog into memory.")
	var syslogHost = flag.String("syslog-host", "0.0.0.0", "Syslog host to listen on.")
	var syslogPort = flag.Int("syslog-port", 10514, "Syslog port to listen on.")
	var host = flag.String("host", "0.0.0.0", "Service host to listen on.")
	var port = flag.Int("port", 3000, "Service port to listen on.")
	var isParseSev = flag.Bool("sev", false, "Parse the syslog severity header")
	flag.Parse()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	sysLogServerDetail := fmt.Sprintf("%s:%d", *syslogHost, *syslogPort)
	server.ListenUDP(sysLogServerDetail)
	server.ListenTCP(sysLogServerDetail)
	server.Boot()

	logArr := make([]string, *maxLogEntries, *maxLogEntries)
	var writeIdx int
	var readIdx int
	broker := sse.NewServer()
	fmt.Printf("Syslog collector started on: %s \n", sysLogServerDetail)

	go func(channel syslog.LogPartsChannel) {
		var logEntry string
		for logParts := range channel {
			// fmt.Println(logParts)
			logEntry = *parseLogEntry(logParts, *isParseSev)
			newWriteIdx := writeIdx + 1
			if newWriteIdx >= *maxLogEntries {
				newWriteIdx = 0
			}
			logArr[newWriteIdx] = logEntry
			writeIdx = newWriteIdx
			// fmt.Printf(logArr[newWriteIdx])
		}
	}(channel)

	ticker := time.NewTicker(time.Duration(*logReadInteval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for readIdx != writeIdx {
					fmt.Printf(logArr[readIdx])
					broker.Notifier <- []byte(logArr[readIdx])
					readIdx++
					if readIdx == *maxLogEntries {
						readIdx = 0
					}
				}

			}
		}
	}()

	serverDetail := fmt.Sprintf("%s:%d", *host, *port)
	log.Fatal("HTTP server error: ", http.ListenAndServe(serverDetail, broker))
	server.Wait()
}

func parseLogEntry(logParts format.LogParts, isParseSev bool) *string {
	// RFC3164
	// 	"timestamp": p.header.timestamp,
	// 	"hostname":  p.header.hostname,
	// 	"tag":       p.message.tag,
	// 	"content":   p.message.content,
	// 	"priority":  p.priority.P,
	// 	"facility":  p.priority.F.Value,
	// 	"severity":  p.priority.S.Value,

	// RFC5424
	// "priority":        p.header.priority.P,
	// "facility":        p.header.priority.F.Value,
	// "severity":        p.header.priority.S.Value,
	// "version":         p.header.version,
	// "timestamp":       p.header.timestamp,
	// "hostname":        p.header.hostname,
	// "app_name":        p.header.appName,
	// "proc_id":         p.header.procId,
	// "msg_id":          p.header.msgId,
	// "structured_data": p.structuredData,
	// "message":         p.message,

	ts := logParts["timestamp"]
	hostname := logParts["hostname"]
	tag := logParts["tag"]
	if tag == nil {
		tag = logParts["app_name"]
	}
	msg := logParts["message"]
	if msg == nil {
		msg = logParts["content"]
	}
	var logStr string
	if isParseSev {
		sev := parseSeverity(logParts["severity"])
		logStr = fmt.Sprintf("[%s][%s][%s][%s]: %s\n", ts, hostname, tag, sev, msg)
	} else {
		logStr = fmt.Sprintf("[%s][%s][%s]: %s\n", ts, hostname, tag, msg)
	}
	// logStr := fmt.Sprintf("[%s][%s][%s][%s]: %s\n", ts, hostname, tag, sev, msg)
	return &logStr
}

func parseSeverity(sev interface{}) string {
	sevNum, ok := sev.(int)
	if !ok {
		return ""
	}
	switch sevNum {
	case 0:
		return "emerg"
	case 1:
		return "alert"
	case 2:
		return "crit"
	case 3:
		return "err"
	case 4:
		return "warning"
	case 5:
		return "notice"
	case 6:
		return "info"
	case 7:
		return "debug"
	}
	return ""
}
