package main

import (
	"dns-tunnel/server/server/session"
	"encoding/base32"
	"fmt"
	"log"
	"os"
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

		fileName := fmt.Sprintf("received%d.txt", currentSession.SessionID)
		err = os.WriteFile(fileName, data, 0644)

		if err != nil {
			log.Println(err)
			return
		}

		log.Printf(
			"Reassembled session %d: wrote %d bytes to %s",
			currentSession.SessionID,
			len(data),
			fileName,
		)

		delete(s.sessions, currentSession.SessionID)
	}
}

func (s *Server) ReassembleSession(sessionID int) ([]byte, error) {
	var encodedData strings.Builder
	targetSession := s.sessions[sessionID]

	for i := 0; i < targetSession.TotalChunks; i++ {
		encodedData.WriteString(targetSession.EncodedData[i])
	}

	decodedData, err := decodeBase32(encodedData.String())

	if err != nil {
		return []byte{}, err
	}

	return decodedData, nil
}

func decodeBase32(encodedString string) ([]byte, error) {
	decodedBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(encodedString)

	if err != nil {
		return []byte{}, err
	}

	return decodedBytes, nil
}
