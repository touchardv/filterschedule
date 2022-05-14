package filterschedule

import (
	"net"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func newResponseMessage(req request.Request) *dns.Msg {
	m := new(dns.Msg)

	switch req.QType() {
	case dns.TypeA:
		a := new(dns.A)
		a.A = net.IPv4zero
		a.Hdr = dns.RR_Header{
			Name:   req.QName(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    300,
		}
		m.Answer = []dns.RR{a}

	case dns.TypeAAAA:
		a := new(dns.AAAA)
		a.AAAA = net.IPv6zero
		a.Hdr = dns.RR_Header{
			Name:   req.QName(),
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    300,
		}
		m.Answer = []dns.RR{a}
	}

	m.Authoritative = true
	m.RecursionAvailable = true
	m.Compress = true
	return m
}
