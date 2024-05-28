package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"social-network/backend/pkg/helper"

	"github.com/pkg/errors"
)

type Group struct {
	ID           int
	Admin_id     int
	Name         string
	Descriptions string
	Mentions     []int
	IsRequested  bool
}

func CreateGroup(title, userId, desc string) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	adminId, _ := helper.MakeInt(userId)

	query, err := db.Prepare("INSERT INTO GROUPS(admin_id, group_name, descriptions, isrequested) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "failed from func AddGroups1")
	}

	defer query.Close()
	if _, err := query.Exec(adminId, title, desc, false); err != nil {
		return errors.Wrap(err, "failed from func AddGroups2")
	}

	return nil
}

func GetGroupByID(id string) (*Group, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	idConv, _ := helper.MakeInt(id)

	fmt.Println("ID: ", id)

	group := Group{}
	row := db.QueryRow("SELECT ID, admin_id, group_name, descriptions, mentioned_users, isrequested FROM GROUPS WHERE ID = ?", idConv)
	var tagged *string

	err = row.Scan(&group.ID, &group.Admin_id, &group.Name, &group.Descriptions, &tagged, &group.IsRequested)

	fmt.Println("gID", group.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}
	members, ok, err := GetMentionedId(group.ID, "ID", "GROUPS")
	if !ok {
		return nil, errors.Wrap(err, "failed from func GetGroupByID3")
	}
	group.Mentions = members

	return &group, nil
}

func GetLastGroup() (*Group, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	group := Group{}
	row := db.QueryRow("SELECT ID, admin_id, group_name, descriptions, mentioned_users, isrequested FROM GROUPS ORDER BY ID DESC LIMIT 1")
	var tagged *string

	err = row.Scan(&group.ID, &group.Admin_id, &group.Name, &group.Descriptions, &tagged, &group.IsRequested)

	fmt.Println("gID", group.ID)
	if err != nil {
		return nil, errors.Wrap(err, "error scanning user data")
	}
	members, ok, err := GetMentionedId(group.ID, "ID", "GROUPS")
	if !ok {
		return nil, errors.Wrap(err, "failed from func GetGroups3")
	}
	group.Mentions = members

	return &group, nil
}

func GetGroups() ([]*Group, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rowF, err := db.Query("SELECT ID, admin_id, group_name, descriptions, isrequested FROM GROUPS")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Mentioned users not found
		}
		log.Println("error getting groups: ", err.Error())
		return nil, errors.Wrap(err, "failed from func GetGroups1")
	}

	var groups []*Group
	defer rowF.Close()
	for rowF.Next() {
		group := &Group{}

		err := rowF.Scan(&group.ID, &group.Admin_id, &group.Name, &group.Descriptions, &group.IsRequested)
		if err != nil {
			return nil, errors.Wrap(err, "failed from func GetGroups2")
		}
		members, ok, err := GetMentionedId(group.ID, "ID", "GROUPS")
		if !ok {
			return nil, errors.Wrap(err, "failed from func GetGroups3")
		}
		group.Mentions = members
		groups = append(groups, group)
	}
	return groups, nil
}

func UpdateMentionedUsers(id, mention string) error {
	fmt.Println(id)
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	groupID, _ := helper.MakeInt(id)

	fmt.Println("GROUPE ID: ", groupID)

	var currentMentionedUsers sql.NullString
	err = db.QueryRow("SELECT mentioned_users FROM GROUPS WHERE ID = ?", groupID).Scan(&currentMentionedUsers)
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrap(err, "failed from func UpdateMentionedUsers")
		}
	}

	// Vérifiez si la valeur est NULL
	if currentMentionedUsers.Valid {
		// La valeur n'est pas NULL, alors analysez-la comme une chaîne JSON
		var mentionedUsers []string
		err = json.Unmarshal([]byte(currentMentionedUsers.String), &mentionedUsers)
		if err != nil {
			log.Fatal(err)
		}

		// Ajoutez le nouvel ID au tableau
		mentionedUsers = append(mentionedUsers, mention)

		// Convertissez le tableau en une chaîne JSON
		newMentionedUsers, err := json.Marshal(mentionedUsers)
		if err != nil {
			log.Fatal(err)
		}

		// Mettez à jour la colonne mentioned_users dans la base de données
		_, err = db.Exec("UPDATE GROUPS SET mentioned_users = ? WHERE ID = ?", string(newMentionedUsers), groupID)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// La valeur est NULL, alors créez un nouveau tableau avec le nouvel ID
		mentionedUsers := []string{mention}
		newMentionedUsers, err := json.Marshal(mentionedUsers)
		if err != nil {
			log.Fatal(err)
		}

		// Mettez à jour la colonne mentioned_users dans la base de données avec le nouveau tableau JSON
		_, err = db.Exec("UPDATE GROUPS SET mentioned_users = ? WHERE ID = ?", string(newMentionedUsers), groupID)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("UPDATED GROUP SUCCESSFULLY")

	return nil
}

func UpdateRequested(id string, value bool) error {
	fmt.Println(id)
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return err
	}
	defer db.Close()

	groupID, _ := helper.MakeInt(id)

	fmt.Println("GROUPE ID: ", groupID)

	// Mettez à jour la colonne mentioned_users dans la base de données avec le nouveau tableau JSON
	_, err = db.Exec("UPDATE GROUPS SET isrequested = ? WHERE ID = ?", value, groupID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UPDATED REQUESTED SUCCESSFULLY")

	return nil
}

func GetAdmin(id string) (int, error) {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		return 0, errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	groupID, _ := helper.MakeInt(id)

	var adminID int
	row := db.QueryRow("SELECT admin_id FROM GROUPS WHERE ID= ?", groupID)

	err = row.Scan(&adminID)

	if err != nil {
		return 0, errors.Wrap(err, "error scanning group chat data")
	}

	return adminID, nil
}