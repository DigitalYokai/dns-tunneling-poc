package session

type Session struct {
	SessionID   int
	TotalChunks int
	EncodedData map[int]string
}

func NewSession(sessionId int, total int) Session {
	return Session{
		SessionID:   sessionId,
		TotalChunks: total,
		EncodedData: make(map[int]string),
	}
}
