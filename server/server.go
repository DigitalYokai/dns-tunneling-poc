package main

import "dns-tunnel/server/server/session"

type Server struct {
	sessions map[int]session.Session
}

func (s *Server) handleChunk(newChunk Chunk) {

}
