package handler

type HandlerNames uint32

const (
	PeerWebsocketHandler HandlerNames = iota
	GetPeerIDHandler
	UserAccountLoginHandler
	PostArticle
	SearchArticle
)
