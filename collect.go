package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/icha024/go-collect-logs/sse"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	var maxLogEntries = flag.Int("max-log-entries", 50000, "Maximum number of log entries to keep. Approx 1KB/entry.")
	var maxFilterEntries = flag.Int("max-filter-entries", 100, "Maximum number of fitlered log entries to return.")
	var logReadInteval = flag.Int("log-read-inteval", 3, "Interval, in seconds, to read syslog into memory.")
	var syslogHost = flag.String("syslog-host", "0.0.0.0", "Syslog host to listen on.")
	var syslogPort = flag.Int("syslog-port", 10514, "Syslog port to listen on.")
	var host = flag.String("host", "0.0.0.0", "Service host to listen on.")
	var port = flag.Int("port", 3000, "Service port to listen on.")
	var enableParseSev = flag.Bool("sev", false, "Parse the syslog severity header")
	var enableStdout = flag.Bool("stdout", true, "Print syslog received to stdout")
	flag.Parse()

	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	syslogServerDetail := fmt.Sprintf("%s:%d", *syslogHost, *syslogPort)
	server.ListenUDP(syslogServerDetail)
	server.ListenTCP(syslogServerDetail)
	server.Boot()

	logArr := make([]string, *maxLogEntries, *maxLogEntries)
	var writeIdx int
	broker := sse.NewServer()
	fmt.Printf("Syslog collector started on: %s \n", syslogServerDetail)

	go func(channel syslog.LogPartsChannel) {
		var logEntry string
		for logParts := range channel {
			// fmt.Println(logParts)
			logEntry = *parseLogEntry(logParts, *enableParseSev)
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
		var readIdx int
		for {
			select {
			case <-ticker.C:
				var buf bytes.Buffer // A Buffer needs no initialization.
				// for readIdx != writeIdx {
				// 	fmt.Printf(logArr[readIdx])
				// 	// broker.Notifier <- []byte(logArr[readIdx])
				// 	buf.Write([]byte(logArr[readIdx]))
				// 	readIdx++
				// 	if readIdx == *maxLogEntries {
				// 		readIdx = 0
				// 	}
				// }

				// if *enableStdout {
				// 	for readIdx != writeIdx {
				// 		fmt.Printf(logArr[readIdx])
				// 		readIdx++
				// 		if readIdx == *maxLogEntries {
				// 			readIdx = 0
				// 		}
				// 	}
				// }

				tmp := writeIdx
				searchIdx := tmp
				for readIdx != searchIdx {
					// fmt.Printf(logArr[searchIdx])
					buf.Write([]byte(logArr[searchIdx]))
					searchIdx--
				}
				if *enableStdout {
					for readIdx != writeIdx {
						fmt.Printf(logArr[readIdx])
						readIdx++
						if readIdx == *maxLogEntries {
							readIdx = 0
						}
					}
				}
				readIdx = tmp
				broker.Notifier <- buf.Bytes()
			}
		}
	}()

	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		query, err := url.QueryUnescape(r.URL.Query().Get("q"))
		if err != nil || len(query) == 0 {
			// log.Println("Invalid query")
			return
		}

		searchIdx := writeIdx
		// if searchIdx < 0 {
		// 	searchIdx = *maxLogEntries - 1
		// }
		for i := 0; i < *maxLogEntries; i++ {
			// fmt.Printf(logArr[readIdx])
			if searchIdx < 0 {
				searchIdx = *maxLogEntries - 1
			}
			logEntry := logArr[searchIdx]
			match := strings.Contains(logEntry, query)
			if match {
				fmt.Fprintf(w, "%s", logArr[searchIdx])
			}
			if i > *maxFilterEntries {
				return
			}
			searchIdx--
		}
		// fmt.Fprintf(w, "Query: %s", query)
	})
	http.Handle("/stream", broker)
	serverDetail := fmt.Sprintf("%s:%d", *host, *port)
	log.Fatal("HTTP server error: ", http.ListenAndServe(serverDetail, nil))
	server.Wait()
}

func parseLogEntry(logParts format.LogParts, enableParseSev bool) *string {
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
	if enableParseSev {
		sev := parseSeverity(logParts["severity"])
		logStr = fmt.Sprintf("[%s][%s][%s][%s]: %s\n", ts, hostname, tag, sev, msg)
	} else {
		logStr = fmt.Sprintf("[%s][%s][%s]: %s\n", ts, hostname, tag, msg)
	}
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
