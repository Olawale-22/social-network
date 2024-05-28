package queries

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

type Notification struct {
	ID         int
	Sender     int
	Recipients []int
	Type       string
	ItemID     int
	Time       string
}

func CreateNotification(sender, recipients, notifType, itemID string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	id, _ := helper.MakeInt(itemID)

	query, err := db.Prepare("INSERT INTO NOTIFICATIONS(sender_id, recipients_id, notification_type, item_id, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertNotifications1")
	}

	defer query.Close()
	if _, err := query.Exec(sender, recipients, notifType, id, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertNotifications2")
	}

	return nil
}

func GetLastNotification() (*Notification, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	notification := Notification{}

	var receivers string
	row := db.QueryRow("SELECT * FROM NOTIFICATIONS ORDER BY ID DESC LIMIT 1")
	if err := row.Scan(
		&notification.ID, &notification.Sender, &receivers, &notification.Type, &notification.ItemID, &notification.Time); err != nil {
		return nil, errors.Wrap(err, "error scanning notif data")
	}

	notification.Recipients = helper.SliceAtoi(receivers)

	return &notification, nil
}

func GetNotificationsByReceiver(receiver int) ([]*Notification, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	receiverConv := strconv.Itoa(receiver)

	rowF, err := db.Query("SELECT * FROM NOTIFICATIONS WHERE recipients_id LIKE ?", "%"+receiverConv+"%")
	if err != nil {
		log.Println("error getting Notifications: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetNotifications1")
	}

	var receivers string
	var notifications []*Notification
	defer rowF.Close()
	for rowF.Next() {
		notification := &Notification{}
		err := rowF.Scan(&notification.ID, &notification.Sender, &receivers, &notification.Type, &notification.ItemID, &notification.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetNotificationsByReceiver1")
		}
		notification.Recipients = helper.SliceAtoi(receivers)
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func DeleteCheckedLastNofification(id int) (int, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Println("Error opening database for deleting notification:", err)
		return 0, err
	}
	defer db.Close()

	var notifID int
	row := db.QueryRow("SELECT ID FROM NOTIFICATIONS WHERE ID = ?", id)
	if err := row.Scan(&notifID); err != nil {
		return 0, errors.Wrap(err, "error scanning notif data")
	}

	_, err = db.Exec("DELETE FROM NOTIFICATIONS WHERE ID = ?", id)
	if err != nil {
		log.Println("Error deleting session from database:", err)
		return 0, err
	}

	return id, nil
}
