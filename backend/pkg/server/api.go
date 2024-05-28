package server

import (
	"database/sql"
	"fmt"

	"social-network/backend/pkg/db/sqlite/queries"

	"github.com/pkg/errors"
)

type Api struct {
	UserInfo      *queries.User           `json:"userInfos"`
	AllUsers      []*queries.User         `json:"users"`
	Follows       []int                   `json:"follows"`
	Followers     []int                   `json:"followers"`
	Posts         *queries.Posts          `json:"posts"`
	Comments      []*queries.Comment      `json:"comments"`
	PrivateChats  []*queries.PrivateChat  `json:"privateChats"`
	Groups        []*queries.Group        `json:"groups"`
	GroupChats    []*queries.GroupChat    `json:"groupChats"`
	Events        []*queries.Event        `json:"events"`
	Votes         []*queries.Vote         `json:"votes"`
	Notifications []*queries.Notification `json:"notifications"`
}

func FetchDataFromDB(id int) (*Api, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, email, firstname, lastname, birthdate, avatar, nickname, aboutme, isconnected, ispublic FROM USERS")
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB1")
	}
	defer rows.Close()

	users := &Api{}
	for rows.Next() {
		user := &queries.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.BirthDate, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsConnected, &user.IsPublic)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func FetchDataFromDB2")
		}
		users.AllUsers = append(users.AllUsers, user)
	}

	// get follows
	follows, err := queries.GetFollows(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB3")
	}

	users.Follows = follows
	fmt.Println("GET-FOLLOWS: ", users.Follows)

	followers, err := queries.GetFollowers(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB3")
	}

	users.Followers = followers
	fmt.Println("GET-FOLLOWERS: ", users.Followers)

	// get user groups
	groups, err := queries.GetGroups()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB6 exactly func GetGroups")
	}
	users.Groups = groups
	chats, err := queries.GetDisplayChat()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB4")
	}
	users.GroupChats = chats

	// get Posts
	posts, err := queries.GetDisplayPost(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB4")
	}
	users.Posts = posts

	// get userInfo
	uzer := queries.User{}
	row := db.QueryRow("SELECT * FROM USERS WHERE ID=?", id)
	if err := row.Scan(
		&uzer.ID, &uzer.Email, &uzer.Password, &uzer.FirstName, &uzer.LastName,
		&uzer.BirthDate, &uzer.Avatar, &uzer.Nickname, &uzer.AboutMe, &uzer.IsConnected, &uzer.IsPublic, &uzer.IsRequested,
	); err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB5")
	}
	users.UserInfo = &uzer

	comments, err := queries.GetComments()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB6")
	}
	users.Comments = comments

	privateChats, err := queries.GetPrivateChats()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB7")
	}
	users.PrivateChats = privateChats

	events, err := queries.GetEvents()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB8")
	}
	users.Events = events

	votes, err := queries.GetVotes()
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB8")
	}
	users.Votes = votes

	notifications, err := queries.GetNotificationsByReceiver(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed from func FetchDataFromDB9")
	}
	users.Notifications = notifications

	return users, nil
}
