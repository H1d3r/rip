package dns_server

import (
	"net"
	"strings"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
)

func typeToString(reqType uint16) string {
	if t, ok := dns.TypeToString[reqType]; ok {
		return t
	}
	return "unknown"
}

func defaultIp(reqType uint16) net.IP {
	if reqType == dns.TypeA {
		return cfg.IPv4
	}
	return cfg.IPv6
}

func parseIp(reqType uint16, name string) net.IP {
	if indx := strings.LastIndex(name, "."); indx != -1 {
		name = name[indx+1:]
	}

	dotCounts := strings.Count(name, "-")
	switch reqType {
	case dns.TypeA:
		if dotCounts != 3 {
			return defaultIp(dns.TypeA)
		}
		return net.ParseIP(strings.Replace(name, "-", ".", -1))
	case dns.TypeAAAA:
		if dotCounts < 2 {
			return defaultIp(dns.TypeAAAA)
		}
		return net.ParseIP(strings.Replace(name, "-", ":", -1))
	default:
		return defaultIp(dns.TypeA)
	}
}
