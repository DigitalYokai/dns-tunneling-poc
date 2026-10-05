package main

import (
	"encoding/base32"
	"fmt"

	"github.com/miekg/dns"
)

type Chunk struct {
	ID          int
	TotalChunks int
	Position    int
	Data        string
}

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

func encodeBase32(data string) string {
	encodedString := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(data))

	return encodedString
}

func createChunks(encodedString string, sessionID int, chunkSize int) []Chunk {
	totalChunks := (len(encodedString) + chunkSize - 1) / chunkSize
	chunks := make([]Chunk, 0, totalChunks)

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize

		if end > len(encodedString) {
			end = len(encodedString)
		}

		newChunk := Chunk{
			ID:          sessionID,
			TotalChunks: totalChunks,
			Position:    i,
			Data:        encodedString[start:end],
		}

		chunks = append(chunks, newChunk)
	}

	return chunks
}

func main() {
	dummyText := "hello world"

	fmt.Println(encodeBase32(dummyText))
}
