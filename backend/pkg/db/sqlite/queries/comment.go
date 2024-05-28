package queries

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
)

type Comment struct {
	ID      int
	User_Id int
	Post_Id int
	Content string
	Image   string
	Time    string
}

func InsertComment(id, postId, comment, image string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	query, err := db.Prepare("INSERT INTO COMMENTS(user_id, post_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertComment1")
	}

	defer query.Close()
	if _, err := query.Exec(id, postId, comment, image, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertComment2")
	}

	return nil
}

func GetComments() ([]*Comment, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT user_id, post_id, content, COALESCE(image, '') AS image, created_at FROM COMMENTS")
	if err != nil {
		log.Println("error getting followers: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetFollows1")
	}

	var comments []*Comment
	defer rowF.Close()
	for rowF.Next() {
		comment := &Comment{}
		err := rowF.Scan(&comment.User_Id, &comment.Post_Id, &comment.Content, &comment.Image, &comment.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetFollows2")
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetLastComment(postId string) (*Comment, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	comment := Comment{}
	row := db.QueryRow("SELECT user_id, post_id, content, image, created_at FROM COMMENTS WHERE post_id=? ORDER BY ID DESC LIMIT 1", postId)
	if err := row.Scan(
		&comment.User_Id, &comment.Post_Id, &comment.Content, &comment.Image, &comment.Time); err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}

	return &comment, nil
}
