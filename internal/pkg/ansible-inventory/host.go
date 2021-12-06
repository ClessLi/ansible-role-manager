package ansible_inventory

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

type Host interface {
	//GetIp() net.IPAddr
	Equal(other Host) bool
	Less(other Host) bool
	IsInclude(other Host) bool
	Len() int
	GetIPString() string
}

type host struct {
	ipAddr net.IPAddr
}

//func NewIPv4Host(ip [4]byte) Host {
//	h, _ := NewHost(net.IPv4(ip[0], ip[1], ip[2], ip[3]), "")
//	return h
//}

//func NewHost(ip net.IP, ipv6Zone string) (Host, error) {
//	h := &host{ipAddr: net.IPAddr{
//		IP:   ip,
//		Zone: ipv6Zone,
//	}}
//	return Host(h), nil
//}

func (h host) Equal(other Host) bool {
	switch oHost := other.(type) {
	case host:
		return h.ipAddr.IP.Equal(oHost.ipAddr.IP)
	}
	return false
}

func (h host) Less(other Host) bool {
	if !h.Equal(other) {
		switch oHost := other.(type) {
		case host:
			return isLessIPAddr(h.ipAddr, oHost.ipAddr)
		default:
			return !other.Less(h)
		}
	}
	return false
}

func (h host) IsInclude(other Host) bool {
	return h.Equal(other)
}

func (h host) Len() int {
	return 1
}

func (h host) GetIPString() string {
	return h.ipAddr.IP.String()
}

type hostPatternIPv4 struct {
	pattern ipv4Pattern
}

func newHostPatternIPv4(pattern ipv4Pattern) Host {
	hostPattern := hostPatternIPv4{pattern: pattern}
	return Host(hostPattern)
}

func (h hostPatternIPv4) Equal(other Host) bool {
	switch oHost := other.(type) {
	case hostPatternIPv4:
		return reflect.DeepEqual(h.pattern, oHost.pattern)
	}
	return false
}

func (h hostPatternIPv4) IsInclude(other Host) bool {
	switch oHost := other.(type) {
	case host:
		return h.pattern.IsIncludeIP(oHost.ipAddr.IP)
	case hostPatternIPv4:
		if h.pattern.IsIncludeIP(oHost.pattern.LowIP()) && h.pattern.IsIncludeIP(oHost.pattern.HighIP()) {
			return true
		}
		return false
	}
	return false
}

func (h hostPatternIPv4) Less(other Host) bool {
	if !h.Equal(other) {
		switch oHost := other.(type) {
		case host:
			return h.getHighIPHost().Less(other)
		case hostPatternIPv4:
			return h.getHighIPHost().Less(oHost.getHighIPHost())
		default:
			return !other.Less(h)
		}
	}
	return false
}

func (h hostPatternIPv4) Len() int {
	count := 1
	for i := net.IPv4len - 1; i >= 0; i-- {
		count *= h.pattern[i].Len()
	}
	return count
}

func (h hostPatternIPv4) GetIPString() string {
	var ipString string
	for i := 0; i < net.IPv4len; i++ {
		ipString += h.pattern[i].String()
		if i == net.IPv4len-1 {
			break
		}
		ipString += "."
	}
	return ipString
}

//func (h hostPatternIPv4) getLowIPHost() Host {
//	return host{ipAddr: net.IPAddr{IP: h.pattern.LowIP()}}
//}

func (h hostPatternIPv4) getHighIPHost() Host {
	return host{ipAddr: net.IPAddr{IP: h.pattern.HighIP()}}
}

type ipv4Pattern [4]ipSegment

func newIPv4Pattern(ipv4PatternBytes [4][2]byte) (ipv4Pattern, error) {
	i4p := ipv4Pattern{}
	for i := 0; i < net.IPv4len; i++ {
		ipSeg, err := newIPSegment(ipv4PatternBytes[i][0], ipv4PatternBytes[i][1])
		if err != nil {
			return i4p, fmt.Errorf("build ipv4Pattern failed, cased by: build the %dst ip segment failed, %s", i, err)
		}
		i4p[i] = ipSeg
	}
	return i4p, nil
}

func (i4p ipv4Pattern) LowIP() net.IP {
	return net.IPv4(i4p[0].low, i4p[1].low, i4p[2].low, i4p[3].low)
}

func (i4p ipv4Pattern) HighIP() net.IP {
	return net.IPv4(i4p[0].high, i4p[1].high, i4p[2].high, i4p[3].high)
}

func (i4p ipv4Pattern) IsIncludeIP(ip net.IP) bool {
	for i := net.IPv4len - 1; i >= 0; i-- {
		if ip[i] < i4p[i].low || ip[i] > i4p[i].high {
			return false
		}
	}
	return true
}

type ipSegment struct {
	low  byte
	high byte
}

func newIPSegment(low, high byte) (ipSegment, error) {
	if low > high {
		return ipSegment{
			low:  low,
			high: low,
		}, fmt.Errorf("low ip segment(%d) is largger than the high's(%d)", low, high)
	}
	return ipSegment{
		low:  low,
		high: high,
	}, nil
}

func (i ipSegment) Len() int {
	return int(i.high - i.low + 1)
}

func (i ipSegment) String() string {
	if i.low == i.high {
		return fmt.Sprintf("%d", i.low)
	}
	return fmt.Sprintf("[%d:%d]", i.low, i.high)
}

func isLessIPAddr(x, y net.IPAddr) bool {
	xIPv4 := x.IP.To4()
	yIPv4 := y.IP.To4()
	return isLessIPv4(xIPv4, yIPv4)
}

func isLessIPv4(x, y net.IP) bool {
	for j := 0; j < net.IPv4len; j++ {
		if x[j] == y[j] {
			continue
		}
		if x[j] < y[j] {
			return true
		}
		return false
	}
	return false

}

func ParseHost(s string) Host {
	if strings.ContainsRune(s, '.') && strings.ContainsRune(s, '[') && strings.ContainsRune(s, ':') && strings.ContainsRune(s, ']') {
		return parseHostPatternIPv4(s)
	} else if strings.ContainsRune(s, '.') {
		return parseHostIPv4(s)
	}
	return nil
}

func parseHostIPv4(s string) Host {
	ipv4 := net.ParseIP(s).To4()
	if ipv4 == nil {
		return nil
	}
	h := host{ipAddr: net.IPAddr{IP: ipv4}}
	return Host(h)
}

func parseHostPatternIPv4(s string) Host {
	segStrs := strings.Split(s, ".")
	if len(segStrs) != net.IPv4len {
		return nil
	}
	ipv4PatternBytes := [4][2]byte{}
	for i, segStr := range segStrs {
		var lowByte, highByte byte

		if idx := strings.Index(segStr, ":"); idx > -1 && segStr[0] == '[' && segStr[len(segStr)-1] == ']' {
			low, lok := dtoi(segStr[1:idx])
			high, hok := dtoi(segStr[idx+1 : len(segStr)-1])
			if !lok || !hok {
				return nil
			}

			if low < 0 || low > 254 || high < low || high > 254 {
				return nil
			}

			lowByte = byte(low)
			highByte = byte(high)
		} else {
			low, ok := dtoi(segStr)
			if !ok {
				return nil
			}
			if low < 0 || low > 254 {
				return nil
			}
			lowByte = byte(low)
			highByte = lowByte
		}

		ipv4PatternBytes[i] = [2]byte{lowByte, highByte}
	}
	pattern, err := newIPv4Pattern(ipv4PatternBytes)
	if err != nil {
		return nil
	}
	return newHostPatternIPv4(pattern)
}

const big = 0xFFFFFF

// Decimal to integer.
// Returns number, success.
func dtoi(s string) (n int, ok bool) {
	i := 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, false
		}
	}
	if i == 0 {
		return 0, false
	}
	return n, true
}
