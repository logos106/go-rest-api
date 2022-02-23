package db

import (
	"database/sql"
	"fmt"

	model "github.com/saroopmathur/rest-api/models"
)

func SelectGroupAccessAll() *[]model.GroupAccess2 {
	db := setupDB()

	query := `SELECT id, name
				FROM user_groups
				WHERE status=$1
				ORDER BY id`
	rows, err := db.Query(query, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("SelectGroupAccessAll: %v\n", err)
		return nil
	}

	var gacs []model.GroupAccess2

	for rows.Next() {
		var gid int
		var gname string

		_ = rows.Scan(&gid, &gname)

		// Get the list of allowed ip for each users
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM group_access_control ga INNER JOIN apps a ON ga.app_id=a.id
					WHERE ga.group_id=$1 AND ga.status=$2`
		rows2, err := db.Query(query, gid, STATUS_ACTIVE)
		if err != nil {
			fmt.Printf("SelectGroupAccessAll: %v\n", err)
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

		gacs = append(gacs, model.GroupAccess2{ID: gid, Group: gname, Apps: apps})
	}

	// close database
	// defer db.Close()

	return &gacs
}

func SelectGroupAccess(did int, gname string, gid int) *[]model.App {
	db := setupDB()

	var rows *sql.Rows
	var err error

	// Get the list of allowed ip for each users
	if gid == 0 {
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM group_access_control ga INNER JOIN apps a ON ga.app_id=a.id
					WHERE ga.user_id=(SELECT id FROM user_groups WHERE name=$1 AND domain_id=$2) AND ga.status=$3`
		rows, err = db.Query(query, gname, did, STATUS_ACTIVE)
	} else {
		query := `SELECT a.id, a.name, a.service_id, a.allowed_ips AS allowed
					FROM group_access_control ga INNER JOIN apps a ON ga.app_id=a.id
					WHERE ga.group_id=$1 AND ga.status=$2`
		rows, err = db.Query(query, gid, STATUS_ACTIVE)
	}

	if err != nil {
		fmt.Printf("SelectGroupAccess: %v\n", err)
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

func InsertGac(domainId int, groupName string, groupId int, appName string, appId int) (*model.GroupAccess, error) {
	db := setupDB()

	var query string
	var err error

	var lastInsertID int
	if groupId > 0 && appId > 0 {
		query = `DELETE FROM group_access_control WHERE group_id=$1 AND app_id=$2`
		db.Exec(query, groupId, appId)

		query = `INSERT INTO group_access_control (group_id, app_id, status)
						VALUES ($1, $2, $3) returning id`
		err = db.QueryRow(query, groupId, appId, STATUS_ACTIVE).Scan(&lastInsertID)
	} else if groupId > 0 {
		query = `DELETE FROM group_access_control WHERE group_id=$1 AND app_id=(SELECT id FROM apps WHERE name=$2)`
		db.Exec(query, groupId, appName)

		query = `INSERT INTO group_access_control (group_id, app_id, status)
						VALUES ($1, (SELECT id FROM apps WHERE name=$2), $3) returning id`
		err = db.QueryRow(query, groupId, appName, STATUS_ACTIVE).Scan(&lastInsertID)
	} else if appId > 0 {
		query = `DELETE FROM group_access_control WHERE group_id=(SELECT id FROM user_groups WHERE domain_id=$1 AND name=$2) AND app_id=$3`
		db.Exec(query, domainId, groupName, appId)

		query = `INSERT INTO group_access_control (group_id, app_id, status)
						VALUES ((SELECT id FROM user_groups WHERE domain_id=$1 AND name=$2), $3, $4) returning id`
		err = db.QueryRow(query, domainId, groupName, appId, STATUS_ACTIVE).Scan(&lastInsertID)
	} else {
		query = `DELETE FROM group_access_control WHERE group_id=(SELECT id FROM user_groups WHERE domain_id=$1 AND name=$2) AND app_id=(SELECT id FROM apps WHERE name=$3)`
		db.Exec(query, domainId, groupName, appName)

		query = `INSERT INTO group_access_control (group_id, app_id, status)
						VALUES ((SELECT id FROM user_groups WHERE domain_id=$1 AND name=$2), (SELECT id FROM apps WHERE name=$3), $4) returning id`
		err = db.QueryRow(query, domainId, groupName, appName, STATUS_ACTIVE).Scan(&lastInsertID)
	}

	// Select the inserted record and return
	return SelectGac(lastInsertID), err
}

func SelectGacs() []model.GroupAccess {
	db := setupDB()

	query := `SELECT id, group_id, allowed_ip_id FROM group_access_control
						WHERE status=$1
						ORDER BY id`
	rows, err := db.Query(query, STATUS_ACTIVE)
	checkErr(err)

	var gacs []model.GroupAccess

	for rows.Next() {
		var id int
		var group int
		var allowed int

		err = rows.Scan(&id, &group, &allowed)
		checkErr(err)
		gacs = append(gacs, model.GroupAccess{ID: id, Group: group})
	}

	// close database
	// defer db.Close()

	return gacs
}

func SelectGac(gid int) *model.GroupAccess {
	db := setupDB()

	query := `SELECT id, group_id, app_id FROM group_access_control
						WHERE status='A' AND id=$1`
	rows, err := db.Query(query, gid)
	checkErr(err)

	var gac *model.GroupAccess

	var id int
	var group int
	var allowed int

	if rows.Next() {
		_ = rows.Scan(&id, &group, &allowed)
		gac = &model.GroupAccess{ID: id, Group: group, App: allowed}
	}

	return gac
}

func DeleteGac(domainId int, groupName string, groupId int, appName string, appId int) int {
	db := setupDB()

	var query string
	var result sql.Result

	if groupId > 0 && appId > 0 {
		query = `UPDATE group_access_control SET status=$1 WHERE group_id=$2 AND app_id=$3`
		result, _ = db.Exec(query, STATUS_DELETED, groupId, appId)
	} else if groupId > 0 {
		query = `UPDATE group_access_control SET status=$1 WHERE group_id=$2 AND app_id=(SELECT id FROM apps WHERE name=$3)`
		result, _ = db.Exec(query, STATUS_DELETED, groupId, appName)
	} else if appId > 0 {
		query = `UPDATE group_access_control SET status=$1 WHERE group_id=(SELECT id FROM user_groups WHERE name=$2 AND domain_id=$3) AND app_id=$4`
		result, _ = db.Exec(query, STATUS_DELETED, groupName, domainId, appId)
	} else {
		query = `UPDATE group_access_control SET status=$1 WHERE group_id=(SELECT id FROM user_groups WHERE name=$2 AND domain_id=$3) AND app_id=(SELECT id FROM apps WHERE name=$4)`
		result, _ = db.Exec(query, STATUS_DELETED, groupName, domainId, appName)
	}

	rowsAffected, _ := result.RowsAffected()

	return int(rowsAffected)
}
