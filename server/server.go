package main

import (
	"dns-tunnel/server/server/session"
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
		data := s.ReassembleSession(currentSession.SessionID)
		log.Printf("Reassembled data: %s", data)

		delete(s.sessions, currentSession.SessionID)
	}
}

func (s *Server) ReassembleSession(sessionID int) string {
	var data strings.Builder
	targetSession := s.sessions[sessionID]

	for i := 0; i < targetSession.TotalChunks; i++ {
		data.WriteString(targetSession.EncodedData[i])
	}

	return data.String()
}
