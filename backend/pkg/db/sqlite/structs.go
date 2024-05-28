package sqlite

import "social-network/backend/pkg/db/sqlite/queries"

type Connection struct {
	Id          int
	IsConnected bool
}

type ServToClientMessage struct {
	Event string
	Data  interface{}
}

type Cchats struct {
	Chats []*queries.GroupChat
}
