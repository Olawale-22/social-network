package queries

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

type GroupChat struct {
	ID       int
	SenderID int
	GroupID int
	Content  string
	Time     string
}

// *********$$$$$$$$$$$$$$$******$
func InsertGroupChat(sender_id, group_id, content string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	sID, _ := helper.MakeInt(sender_id)
	gID, _ := helper.MakeInt(group_id)

	query, err := db.Prepare("INSERT INTO GROUPCHATS(user_id, group_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertPosts1")
	}

	defer query.Close()
	if _, err := query.Exec(sID, gID, content, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertPosts2")
	}
	return nil
}

// *********$$$$$$$$$$$$$$$******$
func GetLastGroupChat() (*GroupChat, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	gChat := GroupChat{}
	row := db.QueryRow("SELECT ID, user_id, group_id, content, created_at FROM GROUPCHATS ORDER BY ID DESC LIMIT 1")

	err = row.Scan(&gChat.ID, &gChat.SenderID, &gChat.GroupID, &gChat.Content, &gChat.Time)

	fmt.Println("gID", gChat.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}

	return &gChat, nil
}

// $$$$$$$$$$*********$$$$$$$$$$$$$$***********
func GetDisplayChat() ([]*GroupChat, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT ID, user_id, group_id, content, created_at FROM GROUPCHATS")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Mentioned users not found
		}
		log.Println("error getting groups: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetGroups1")
	}

	var groupChats []*GroupChat
	defer rowF.Close()
	for rowF.Next() {
		groupChat := &GroupChat{}

		err := rowF.Scan(&groupChat.ID, &groupChat.SenderID, &groupChat.GroupID, &groupChat.Content, &groupChat.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetGroups2")
		}
		groupChats = append(groupChats, groupChat)
	}
	return groupChats, nil
}

func GetGroupChatMembers(id int) (string, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return "", errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	var results string
	row := db.QueryRow("SELECT COALESCE(g.mentioned_users, '') FROM GROUPCHATS gc INNER JOIN GROUPS g ON g.ID = gc.group_id WHERE g.ID = ?", id)
	err = row.Scan(&results)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", errors.Wrap(err, "error scanning user data")
	}

	return results, nil
}