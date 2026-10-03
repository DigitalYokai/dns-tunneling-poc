package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func queryDNS(domain, server string, qtype uint16) ([]string, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.Net = "udp"
	resp, _, err := c.Exchange(m, server+":53")

	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("query error: %s", dns.RcodeToString[resp.Rcode])
	}

	var results []string

	for _, ans := range resp.Answer {
		switch qtype {
		case dns.TypeA:
			if a, ok := ans.(*dns.A); ok {
				results = append(results, a.A.String())
			}
		case dns.TypeAAAA:
			if aaaa, ok := ans.(*dns.AAAA); ok {
				results = append(results, aaaa.AAAA.String())
			}
		}
	}

	return results, nil
}

func handleA(w dns.ResponseWriter, r *dns.Msg) {
	question := r.Question[0]

	//logs the dns query, nned to decode data from it later
	log.Printf("Received query for: %s", question.Name)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	//returns an empty answer for now
	w.WriteMsg(m)
}

func main() {
	dns.HandleFunc(".", handleA)
	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Println("DNS server running on :8053")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
