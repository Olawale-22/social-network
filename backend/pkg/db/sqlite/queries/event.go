package queries

import (
	"database/sql"
	"time"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

type Event struct {
	ID          int
	UserID      int
	GroupID    int
	Title       string
	Description string
	Time        string
	CreatedAt   string
}

func CreateEvent(userID, groupeID, title, description, eventTime string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()

	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	query, err := db.Prepare("INSERT INTO EVENTS(user_id, group_id, title, description, time, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertEvent1")
	}

	defer query.Close()
	if _, err := query.Exec(userID, groupeID, title, description, eventTime, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertEvent2")
	}

	return nil
}

func GetLastEvent(groupeID string) (*Event, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	id, _ := helper.MakeInt(groupeID)
	event := Event{}
	row := db.QueryRow("SELECT * FROM EVENTS WHERE group_id = ? ORDER BY ID DESC LIMIT 1", id)
	if err := row.Scan(
		&event.ID, &event.UserID, &event.GroupID, &event.Title, &event.Description, &event.Time, &event.CreatedAt); err != nil {
		return nil, errors.Wrap(err, "error scanning event data")
	}

	return &event, nil
}

func GetEvents() ([]*Event, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT * FROM EVENTS")
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetEvents1")
	}

	var events []*Event
	defer rowF.Close()
	for rowF.Next() {
		event := &Event{}
		err := rowF.Scan(&event.ID, &event.UserID, &event.GroupID, &event.Title, &event.Description, &event.Time, &event.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetEvents2")
		}
		events = append(events, event)
	}
	return events, nil
}