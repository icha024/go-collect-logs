package main

import (
	"bufio"
	"fmt"
	"gopkg.in/mcuadros/go-syslog.v2"
	"os"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	// server.SetFormat(syslog.RFC5424)
	// server.SetFormat(syslog.RFC3164)
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:10514")
	server.ListenTCP("0.0.0.0:10514")
	server.Boot()
	fmt.Println("Log collector started.")

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
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
			fmt.Printf("[%s][%s][%s] %s\n", logParts["timestamp"], logParts["hostname"], logParts["tag"], logParts["content"])
		}
	}(channel)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		// if text == "quit" {
		// 	break
		// }
	}
	server.Wait()
}
