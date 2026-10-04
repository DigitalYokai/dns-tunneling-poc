package session

type Session struct {
	SessionID   int
	TotalChunks int
	encodedData map[int]string
}

func newSession(sessionId int, total int) Session {
	return Session{
		SessionID:   sessionId,
		TotalChunks: total,
		encodedData: make(map[int]string),
	}
}
