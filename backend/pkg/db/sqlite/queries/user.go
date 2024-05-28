package queries

import (
	"database/sql"
	"fmt"
	"log"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	Password    string
	Email       string
	FirstName   string
	LastName    string
	BirthDate   string
	Avatar      string
	Nickname    string
	Followers   []int
	AboutMe     string
	Uuid        string
	IsConnected bool
	IsPublic    bool
	IsRequested bool
}

func CreateUser(email, pass, firstName, lastName, birthDate, avatar, nickname, aboutme string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()
	password := GenerateHash(pass)

	query, err := db.Prepare("INSERT INTO USERS(email, password, firstname, lastname, birthdate, avatar, nickname, aboutme, isconnected, ispublic, isrequested) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer query.Close()
	if _, err := query.Exec(email, password, firstName, lastName, birthDate, avatar, nickname, aboutme, false, true, false); err != nil {
		return err
	}
	return nil
}

func GetByEmailOrNickname(email, nickname string) (bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM USERS WHERE email=? OR nickname=?", email, nickname).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func VerifyLogin(username, pass string) (bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var pss string
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM USERS WHERE (email=? OR nickname=?)", username, username).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = db.QueryRow("SELECT password FROM USERS WHERE (email=? OR nickname=?)", username, username).Scan(&pss)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return ComparePassword(pss, pass), nil
}

func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetInfo(u string) (*User, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	_, err = db.Exec("UPDATE USERS SET isconnected=? WHERE (email=? OR nickname=?)", true, u, u)
	if err != nil {
		return nil, errors.Wrap(err, "error updating isconnected")
	}

	user := User{}
	row := db.QueryRow("SELECT * FROM USERS WHERE (email=? OR nickname=?)", u, u)
	if err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.BirthDate, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsConnected, &user.IsPublic, &user.IsRequested); err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}

	// Fetch follower IDs and populate the Followers array
	rows, err := db.Query("SELECT follower_id FROM FOLLOWERS WHERE user_id=?", user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching follower IDs")
	}
	defer rows.Close()

	for rows.Next() {
		var followerID int
		if err := rows.Scan(&followerID); err != nil {
			return nil, errors.Wrap(err, "error scanning follower ID")
		}
		user.Followers = append(user.Followers, followerID)
	}

	return &user, nil
}

func GenerateHash(password string) string {
	var mpCrypt []byte
	cost := 11
	mpCrypt, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Println("Error in GenerateHash:")
		log.Fatal(err)
	}
	return string(mpCrypt)
}

func GetUsers(s string) ([]*User, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM USERS WHERE nickname != ?", s)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetUsers1")
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsConnected, &user.IsPublic, &user.IsRequested)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetUsers2")
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed from func GetUsers3")
	}

	return users, nil
}

func GetConnUsers(id int) ([]*User, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetConnUsers1")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM USERS WHERE IsConnected = ? AND ID != ?", true, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetConnUsers2")
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsConnected, &user.IsPublic, &user.IsRequested)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetConnUsers3")
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed from func GetConnUsers4")
	}

	return users, nil
}

func GetPrivacyByID(id int) (bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		fmt.Println("Error opening database: ", err.Error())
		return false, err
	}
	defer db.Close()
	var isPublic bool
	err = db.QueryRow("SELECT ispublic FROM USERS WHERE ID=?", id).Scan(&isPublic)
	if err != nil {
		fmt.Println("Error getting privacy: ", err.Error())
		return false, err
	}
	return isPublic, nil
}

func UpdatePrivacy(id int, state bool) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE USERS SET ispublic=? WHERE ID=?", state, id)
	if err != nil {
		return err
	}
	return nil
}

type CustomError string

func (ce CustomError) Error() string {
	return string(ce)
}

func CCustomError() error {
	return CustomError("FollowerExist")
}

// END*****$$$$$$$$$$$*************$$$$$$$$$******$

func GetMentionedId(pid interface{}, columnName, table string) ([]int, bool, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, false, err
	}

	defer db.Close()
	var blob string
	query := "SELECT COALESCE(mentioned_users, '') AS mentioned_users FROM " + table + " WHERE " + columnName + " = ?"
	err = db.QueryRow(query, pid).Scan(&blob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil // Mentioned users not found
		}
		return nil, false, errors.Wrap(err, "failed from func GetMentionedId1")
	}

	ok := helper.SliceAtoi(blob)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed from func GetMentionedId2")
	}
	return ok, true, nil
}

func GetInfoById(id int) (*User, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	user := User{}
	row := db.QueryRow("SELECT * FROM USERS WHERE (ID=?)", id)
	if err := row.Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.BirthDate, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsConnected, &user.IsPublic, &user.IsRequested,
	); err != nil {
		return nil, errors.Wrap(err, "error scanning user data GetInfoByID1")
	}

	// Fetch follower IDs and populate the Followers array
	rows, err := db.Query("SELECT follower_id FROM FOLLOWERS WHERE user_id=?", id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func GetInfoByID2")
	}
	defer rows.Close()

	var followerID int
	for rows.Next() {
		if err := rows.Scan(&followerID); err != nil {
			return nil, errors.Wrap(err, "failed from func GetInfoByID3")
		}
		user.Followers = append(user.Followers, followerID)
	}

	return &user, nil
}

func UpdateRequestedUser(id string, value bool) error {
	fmt.Println(id)
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	groupID, _ := helper.MakeInt(id)

	fmt.Println("GROUPE ID: ", groupID)

	// Mettez à jour la colonne mentioned_users dans la base de données avec le nouveau tableau JSON
	_, err = db.Exec("UPDATE USERS SET isrequested = ? WHERE ID = ?", value, groupID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UPDATED REQUESTED USERS SUCCESSFULLY")

	return nil
}
