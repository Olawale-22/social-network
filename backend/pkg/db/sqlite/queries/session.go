package queries

import (
	"database/sql"
	"fmt"
	"log"
)

func AddSession(name string, uuid string) error {

	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Println("error opening db for adding session")
		return err
	}

	defer db.Close()

	if err != nil {
		fmt.Println("erreur creation de table session")
		return err
	}

	if name != "" && uuid != "" {

		//verifier si l'utilisateur est déjà connecté
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM SESSION  WHERE nickname = ?", name).Scan(&count)
		if err != nil {
			log.Println("Erreur lors de la recherche de l'utilisateur dans la base de données, func AddSession:")
			return err

		}
		if count == 0 {

			_, err := db.Exec("INSERT INTO SESSION (nickname, uuid) VALUES (?, ?)", name, uuid)
			if err != nil {
				log.Println("Erreur à l'insertion de donnée dans session, func AddSession:")
				return err
			}

			_, err = db.Exec("UPDATE USERS SET isconnected = ? WHERE nickname = ?", true, name)
			if err != nil {
				log.Println("Erreur Update isconnected top USERS, func AddSession:")
				return err
			}

			// _, err = db.Exec("UPDATE USERS SET uuid = ? WHERE nickname = ?", uuid, name)
			// if err != nil {
			// 	fmt.Println("Erreur Update uuid top USERS, func AddSession:")
			// 	log.Fatal(err)
			// }
		}
	} else {
		fmt.Println("name uuid cookie vide !")
	}
	return nil
}

/*************************** DELETE SESSION **********************************/
func DeleteSession(id int) error {
	db, err := sql.Open("sqlite3", "backend/pkg/db/sqlite/sqlite.db")
	if err != nil {
		log.Println("Error opening database for deleting session:", err)
		return err
	}
	defer db.Close()

	// Perform the DELETE and UPDATE operations separately
	_, err = db.Exec("DELETE FROM SESSION WHERE ID = ?", id)
	if err != nil {
		log.Println("Error deleting session from database:", err)
		return err
	}

	// Update the user's "isconnected" field
	_, err = db.Exec("UPDATE USERS SET isconnected = ? WHERE ID = ?", false, id)
	if err != nil {
		log.Println("Error updating user's 'isconnected' field:", err)
		return err
	}

	return nil
}
