package trace

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	TraceID = "trace-id"
)

// Trace ...
type Trace struct {
	TraceID string
}

// GetTrace ...
func GetTrace(req *http.Request) *Trace {
	t := &Trace{
		TraceID: req.Header.Get(TraceID),
	}

	if t.TraceID == "" {
		t.TraceID = genTraceID()
		req.Header.Set(TraceID, t.TraceID) // 防止多次生成
	}

	return t
}

func genTraceID() string {
	ip := getLocalIP()
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()
	b := bytes.Buffer{}

	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	b.WriteString(fmt.Sprintf("%x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0")

	return b.String()
}

func genSpanID() string {
	return fmt.Sprintf("%x", rand.Int63())
}

func getLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}
