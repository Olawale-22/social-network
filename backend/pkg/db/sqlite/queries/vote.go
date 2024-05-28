package queries

import (
	"database/sql"
	"log"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

type Vote struct {
	ID         int
	EventID    int
	UserID     int
	VoteOption string
}

func CreateVote(eventID, userID, voteOption string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Prepare("INSERT INTO VOTES(event_id, user_id, vote_option) VALUES (?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertVote1")
	}

	defer query.Close()
	if _, err := query.Exec(eventID, userID, voteOption); err != nil {
		return errors.Wrap(err, "failed from func InsertVote2")
	}

	return nil
}

func GetCountVotes(event, user string) (int, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	eventID, _ := helper.MakeInt(event)
	userID, _ := helper.MakeInt(user)

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM VOTES WHERE event_id = ? AND user_id = ?", eventID, userID)
	if err := row.Scan(&count); err != nil {
		return 0, errors.Wrap(err, "error getting count vote")
	}

	return count, nil
}

func UpdateVote(event, user, newOption string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE VOTES SET vote_option = ? WHERE event_id = ? AND user_id = ?", newOption, event, user)
	if err != nil {
		log.Println("Error upadating vote from database:", err)
		return err
	}

	return nil
}

func GetLastVote(id, user string) (*Vote, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	eventID, _ := helper.MakeInt(id)
	userID, _ := helper.MakeInt(user)

	vote := Vote{}
	row := db.QueryRow("SELECT * FROM VOTES WHERE event_id = ? AND user_id= ? ORDER BY ID DESC LIMIT 1", eventID, userID)
	if err := row.Scan(
		&vote.ID, &vote.EventID, &vote.UserID, &vote.VoteOption); err != nil {
		return nil, errors.Wrap(err, "error scanning event data")
	}

	return &vote, nil
}

func GetVotes() ([]*Vote, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT * FROM VOTES")
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetVotes1")
	}

	var votes []*Vote
	defer rowF.Close()
	for rowF.Next() {
		vote := &Vote{}
		err := rowF.Scan(&vote.ID, &vote.EventID, &vote.UserID, &vote.VoteOption)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetVotes2")
		}
		votes = append(votes, vote)
	}
	return votes, nil
}