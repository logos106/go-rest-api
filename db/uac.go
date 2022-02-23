package db

import (
	"database/sql"
	"fmt"

	model "github.com/saroopmathur/rest-api/models"
)

// Select returns the whole database
func SelectUserAccessAll() *[]model.UserAccess2 {
	db := setupDB()

	query := `SELECT id, name
				FROM users
				WHERE status=$1
				ORDER BY id`
	rows, err := db.Query(query, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("SelectUserAccessAll: %v\n", err)
		return nil
	}

	var uacs []model.UserAccess2

	for rows.Next() {
		var uid int
		var uname string

		_ = rows.Scan(&uid, &uname)

		// Get the list of allowed ip for each users
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM user_access_control ua INNER JOIN apps a ON ua.app_id=a.id
					WHERE ua.user_id=$1 AND ua.status=$2`
		rows2, err := db.Query(query, uid, STATUS_ACTIVE)
		if err != nil {
			fmt.Printf("SelectUserAccessAll: %v\n", err)
			return nil
		}

		var apps []model.App
		for rows2.Next() {
			var aid int
			var aname string
			var serviceId int
			var allowed string

			_ = rows2.Scan(&aid, &aname, &serviceId, &allowed)

			apps = append(apps, model.App{ID: aid, Name: aname, ServiceId: serviceId, AllowedIPs: allowed})
		}

		uacs = append(uacs, model.UserAccess2{ID: uid, User: uname, Apps: apps})
	}

	// close database
	// defer db.Close()

	return &uacs
}

func SelectUserAccess(domainId int, uname string, uid int) *[]model.App {
	db := setupDB()

	var rows *sql.Rows
	var err error

	// Get the list of allowed ip for each users
	if uid == 0 {
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM user_access_control ua INNER JOIN apps a ON ua.app_id=a.id
					WHERE ua.user_id=(SELECT id FROM users WHERE name=$1 AND domain_id=$2) AND ua.status=$3`
		rows, err = db.Query(query, uname, STATUS_ACTIVE)
	} else {
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM user_access_control ua INNER JOIN apps a ON ua.app_id=a.id
					WHERE ua.user_id=$1 AND ua.status=$2`
		rows, err = db.Query(query, uid, STATUS_ACTIVE)
	}

	if err != nil {
		fmt.Printf("SelectUserAccess: %v\n", err)
		return nil
	}

	var apps []model.App
	for rows.Next() {
		var aid int
		var aname string
		var serviceId int
		var allowed string

		_ = rows.Scan(&aid, &aname, &serviceId, &allowed)

		apps = append(apps, model.App{ID: aid, Name: aname, ServiceId: serviceId, AllowedIPs: allowed})
	}

	// close database
	// defer db.Close()

	return &apps
}

// Insert allows populating database
func InsertUac(domainId int, userName string, userId int, appName string, appId int) (*model.UserAccess, error) {
	db := setupDB()

	var query string
	var err error

	var lastInsertID int
	if userId > 0 && appId > 0 {
		query = `DELETE FROM user_access_control WHERE user_id=$1 AND app_id=$2`
		db.Exec(query, userId, appId)

		query = `INSERT INTO user_access_control (user_id, app_id, status)
					VALUES ($1, $2, $3) returning id`
		err = db.QueryRow(query, userId, appId, STATUS_ACTIVE).Scan(&lastInsertID)
	} else if userId > 0 {
		query = `DELETE FROM user_access_control WHERE user_id=$1 AND app_id=(SELECT id FROM apps WHERE name=$2)`
		db.Exec(query, userId, appName)

		query = `INSERT INTO user_access_control (user_id, app_id, status)
					VALUES ($1, (SELECT id FROM apps WHERE name=$2), $3) returning id`
		err = db.QueryRow(query, userId, appName, STATUS_ACTIVE).Scan(&lastInsertID)
	} else if appId > 0 {
		query = `DELETE FROM user_access_control WHERE user_id=(SELECT id FROM users WHERE domain_id=$1 AND name=$2) AND app_id=$3`
		db.Exec(query, domainId, userName, appId)

		query = `INSERT INTO user_access_control (user_id, app_id, status)
					VALUES ((SELECT id FROM users WHERE domain_id=$1 AND name=$2), $3, $4) returning id`
		err = db.QueryRow(query, domainId, userName, appId, STATUS_ACTIVE).Scan(&lastInsertID)

	} else {
		query = `DELETE FROM user_access_control WHERE user_id=(SELECT id FROM users WHERE domain_id=$1 AND name=$2) AND app_id=(SELECT id FROM apps WHERE name=$3)`
		db.Exec(query, domainId, userName, appName)

		query = `INSERT INTO user_access_control (user_id, app_id, status)
					VALUES ((SELECT id FROM users WHERE domain_id=$1 AND name=$2), (SELECT id FROM apps WHERE name=$3), $4) returning id`
		err = db.QueryRow(query, domainId, userName, appName, STATUS_ACTIVE).Scan(&lastInsertID)
	}

	// Select the inserted record and return
	return SelectUac(lastInsertID), err
}

// Delete the record with the id
func DeleteUac(domainId int, userName string, userId int, appName string, appId int) int {
	db := setupDB()

	var query string
	var result sql.Result

	if userId > 0 && appId > 0 {
		query = `UPDATE user_access_control SET status=$1 WHERE user_id=$2 AND app_id=$3`
		result, _ = db.Exec(query, STATUS_DELETED, userId, appId)
	} else if userId > 0 {
		query = `UPDATE user_access_control SET status=$1 WHERE user_id=$2 AND app_id=(SELECT id FROM apps WHERE name=$3)`
		result, _ = db.Exec(query, STATUS_DELETED, userId, appName)
	} else if appId > 0 {
		query = `UPDATE user_access_control SET status=$1 WHERE user_id=(SELECT id FROM users WHERE name=$2 AND domain_id=$3) AND app_id=$4`
		result, _ = db.Exec(query, STATUS_DELETED, userName, domainId, appId)
	} else {
		query = `UPDATE user_access_control SET status=$1 WHERE user_id=(SELECT id FROM users WHERE name=$2 AND domain_id=$3) AND app_id=(SELECT id FROM apps WHERE name=$4)`
		result, _ = db.Exec(query, STATUS_DELETED, userName, domainId, appName)
	}

	rowsAffected, _ := result.RowsAffected()

	return int(rowsAffected)
}

// Selecte the record with the id
func SelectUac(uaid int) *model.UserAccess {
	db := setupDB()

	query := `SELECT id, user_id, app_id
				FROM user_access_control
				WHERE id=$1 AND status=$2`
	rows, err := db.Query(query, uaid, STATUS_ACTIVE)
	checkErr(err)

	var uac *model.UserAccess

	var id int
	var user int
	var app int

	if rows.Next() {
		_ = rows.Scan(&id, &user, &app)
		uac = &model.UserAccess{ID: id, User: user, App: app}
	}

	return uac
}
