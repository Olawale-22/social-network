package queries

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
)

type Mentions struct {
	User_Id  int
	Post_Id  int
	Mentions []int
}

type Post struct {
	PiD      int
	User_Id  int
	Content  string
	Image    string
	Time     string
	Privacy  string
	Mentions []int
}

type Posts struct {
	Public   []*Post
	Private  []*Post
	Mentions []*Post
}

func CreatePost(id, post, privacy, image string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	query, err := db.Prepare("INSERT INTO POSTS(user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertPosts1")
	}

	defer query.Close()
	if _, err := query.Exec(id, post, image, privacy, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertPosts2")
	}

	return nil
}

func GetDisplayPost(hiDee int) (*Posts, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}

	defer db.Close()
	public := "public"
	private := "private"
	mentions := "mentions"

	rowF, err := db.Query("SELECT ID, user_id, content, COALESCE(image, '') AS image, privacy, created_at FROM POSTS WHERE privacy != ?", mentions)
	if err != nil {
		log.Println("error getting followers: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetDisplayPost1")
	}

	defer rowF.Close()
	PostArray := Posts{}
	for rowF.Next() {
		pPost := &Post{}
		err := rowF.Scan(&pPost.PiD, &pPost.User_Id, &pPost.Content, &pPost.Image, &pPost.Privacy, &pPost.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetDisplayPost2")
		}

		if pPost.Privacy == public {
			PostArray.Public = append([]*Post{pPost}, PostArray.Public...)
		} else if pPost.Privacy == private {
			if hiDee == pPost.User_Id {
				PostArray.Private = append(PostArray.Private, pPost)
			} else {
				followers, err := GetFollowers(pPost.User_Id)
				if err != nil {
					return nil, errors.Wrap(err, "failed from func GetDisplayPost3")
				}
				for _, follower := range followers {
					if follower == hiDee {
						PostArray.Private = append(PostArray.Private, pPost)
					}
				}
			}
		}
	}

	rowPr, err := db.Query("SELECT ID, user_id, content, COALESCE(image, '') AS image, privacy, created_at FROM POSTS WHERE privacy=?", mentions)
	if err != nil {
		log.Println("error getting followers: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetDisplayPost4")
	}

	for rowPr.Next() {
		pPost := &Post{}
		err := rowPr.Scan(&pPost.PiD, &pPost.User_Id, &pPost.Content, &pPost.Image, &pPost.Privacy, &pPost.Time)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetDisplayPost5")
		}

		tagged, ok, err := GetMentionedId(pPost.PiD, "post_id", "MENTIONS")
		if !ok {
			if err != nil {
				return nil, errors.Wrap(err, "failed from func GetDisplayPost6")
			}
		} else {
			if pPost.User_Id == hiDee {
				PostArray.Mentions = append(PostArray.Mentions, pPost)
			}
			for _, id := range tagged {
				if id == hiDee {
					PostArray.Mentions = append(PostArray.Mentions, pPost)
				}
			}
		}
	}

	return &PostArray, nil
}

func InsertMentionedPost(id, post, image, privacy, mentions string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}

	defer db.Close()
	datetime := time.Now()
	dt := datetime.Format(time.RFC3339)

	query, err := db.Prepare("INSERT INTO POSTS(user_id, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertMentionedPost1")
	}

	defer query.Close()
	if _, err := query.Exec(id, post, image, privacy, dt); err != nil {
		return errors.Wrap(err, "failed from func InsertMentionedPost2")
	}

	// get post id
	var Pid int
	err = db.QueryRow("SELECT ID FROM POSTS WHERE user_id=? AND content=?", id, post).Scan(&Pid)
	if err != nil {
		return errors.Wrap(err, "failed from func InsertMentionedPost3")
	}

	// insert into mentions here...
	qry, err := db.Prepare("INSERT INTO MENTIONS(user_id, post_id, mentioned_users) VALUES (?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func InsertMentionedPost4")
	}

	defer qry.Close()
	if _, err := qry.Exec(id, Pid, mentions); err != nil {
		return errors.Wrap(err, "failed from func InsertMentionedPost5")
	}
	return nil
}

func GetLastPost() (*Post, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	post := Post{}
	row := db.QueryRow("SELECT ID, user_id, content, image, privacy, created_at FROM POSTS ORDER BY ID DESC LIMIT 1")
	if err := row.Scan(
		&post.PiD, &post.User_Id, &post.Content, &post.Image, &post.Privacy, &post.Time); err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}

	return &post, nil
}
