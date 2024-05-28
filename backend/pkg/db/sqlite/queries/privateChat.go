package queries

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
)

type PrivateChat struct {
	ID        int
	Sender    int
	Recipient int
	Content   string
	Time      string
}

func CreatePrivateChat(senderId, recipient, content string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	query, err := db.Prepare("INSERT INTO PRIVATECHATS(sender_id, recipient_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertMessage1")
	}

	defer query.Close()
	if _, err := query.Exec(senderId, recipient, content, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertMessage2")
	}

	return nil
}

func GetLastPrivateChat() (*PrivateChat, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	conversation := PrivateChat{}
	row := db.QueryRow("SELECT * FROM PRIVATECHATS ORDER BY ID DESC LIMIT 1")
	if err := row.Scan(
		&conversation.ID, &conversation.Sender, &conversation.Recipient, &conversation.Content, &conversation.Time); err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}

	return &conversation, nil
}

func GetMessagesByPrivateChat(senderId, recipients string) ([]*PrivateChat, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	rowF, err := db.Query("SELECT sender_id, recipient_id, content, created_at FROM PRIVATECHATS WHERE sender_id= ? AND recipients LIKE ?", senderId, "%"+recipients+"%")
	if err != nil {
		log.Println("error getting messages: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetMessageByConversation1")
	}

	var conversations []*PrivateChat
	defer rowF.Close()
	for rowF.Next() {
		conversation := &PrivateChat{}
		err := rowF.Scan(&conversation.Sender, &conversation.Recipient, &conversation.Content, dt)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetMessageByConversation2")
		}
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}

func GetPrivateChats() ([]*PrivateChat, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT * FROM PRIVATECHATS")
	if err != nil {
		log.Println("error getting followers: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetFollows1")
	}

	var conversations []*PrivateChat
	defer rowF.Close()
	for rowF.Next() {
		conversation := &PrivateChat{}
		err := rowF.Scan(&conversation.ID, &conversation.Sender, &conversation.Recipient, &conversation.Content, &conversation.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetConversations1")
		}
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}
