package main

import (
	"dns-tunnel/server/server/session"
	"encoding/base32"
	"log"
	"strings"
)

type Server struct {
	sessions map[int]session.Session
}

func (s *Server) handleChunk(newChunk Chunk) {

	currentSession, exists := s.sessions[newChunk.ID]

	// Checks if session already exists, if not creates one
	if !exists {
		currentSession = session.NewSession(
			newChunk.ID,
			newChunk.TotalChunks,
		)
	}

	currentSession.EncodedData[newChunk.Position] = newChunk.Data
	s.sessions[newChunk.ID] = currentSession

	if len(currentSession.EncodedData) == currentSession.TotalChunks {
		// Reassemble and decode
		data, err := s.ReassembleSession(currentSession.SessionID)

		if err != nil {
			log.Printf("Failed to decode session %d: %v",
				currentSession.SessionID,
				err,
			)
			return
		}

		log.Printf("Reassembled data: %s", data)

		delete(s.sessions, currentSession.SessionID)
	}
}

func (s *Server) ReassembleSession(sessionID int) (string, error) {
	var data strings.Builder
	targetSession := s.sessions[sessionID]

	for i := 0; i < targetSession.TotalChunks; i++ {
		data.WriteString(targetSession.EncodedData[i])
	}

	decodedData, err := decodeString32(data.String())

	if err != nil {
		return "", err
	}

	return decodedData, nil
}

func decodeString32(encodedString string) (string, error) {
	decodedBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(encodedString)

	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}
