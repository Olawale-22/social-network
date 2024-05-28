package queries

import (
	"database/sql"
	"log"
	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

func CreateFollower(uID, FlowID string) (bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	userId, _ := helper.MakeInt(uID)
	FollowerId, _ := helper.MakeInt(FlowID)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM FOLLOWERS WHERE user_id = ? AND follower_id = ?", userId, FollowerId).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "failed from func AddFollower1")
	}

	if count > 0 {
		// Relationship already exists
		return false, nil
	}

	query, err := db.Prepare("INSERT INTO FOLLOWERS(user_id, follower_id) VALUES (?, ?)")
	if err != nil {
		return false, errors.Wrap(err, "failed from func AddFollower2")
	}
	defer query.Close()

	if _, err := query.Exec(userId, FollowerId); err != nil {
		return false, errors.Wrap(err, "failed from func AddFollower3")
	}

	return true, nil
}

func GetFollows(id int) ([]int, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT follower_id FROM FOLLOWERS WHERE user_id=?", id)
	if err != nil {
		log.Println("error getting follows: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetFollows1")
	}

	var follows []int
	defer rowF.Close()
	for rowF.Next() {
		var follow int
		err := rowF.Scan(&follow)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetFollows2")
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

func GetFollowers(id int) ([]int, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT user_id FROM FOLLOWERS WHERE follower_id=?", id)
	if err != nil {
		log.Println("error getting followers: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetFollowers1")
	}

	var followers []int
	defer rowF.Close()
	for rowF.Next() {
		var follower int
		err := rowF.Scan(&follower)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetFollowers2")
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

func CheckFollowing(id, followId string) (bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	userId, _ := helper.MakeInt(id)
	FollowerId, _ := helper.MakeInt(followId)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM FOLLOWERS WHERE user_id = ? AND follower_id = ?", userId, FollowerId).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "failed from func CheckFollowing")
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

func Unfollow(id, followId int) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Println("Error opening database for deleting session:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM FOLLOWERS WHERE user_id = ? AND follower_id=?", id, followId)
	if err != nil {
		log.Println("Error deleting follow from database:", err)
		return err
	}

	return nil
}
