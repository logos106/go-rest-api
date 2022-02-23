package db

import (
	"database/sql"
	"fmt"

	model "github.com/saroopmathur/rest-api/models"
)

// Insert allows populating database
func InsertGroup(domainId int, group *model.Group) (*model.Group2, error) {
	db := setupDB()

	var lastInsertID int
	query := `INSERT INTO user_groups (name, domain_id, status)
						VALUES($1, $2, $3) returning id;`
	fmt.Printf("CreateGroup %s domain=%d\n", group.Name, domainId)
	err := db.QueryRow(query, group.Name, domainId, STATUS_ACTIVE).Scan(&lastInsertID)
	if err != nil {
		fmt.Printf("CreateGroup %s domain=%d %v\n", group.Name, domainId, err)
		return nil, err
	}

	// Select the inserted record and return
	return SelectGroup(domainId, group.Name, group.ID), nil
}

// Select returns the whole database
func SelectGroups(domainId int) []*model.Group2 {
	db := setupDB()

	query := `SELECT g.id, g.name, COALESCE(c.cnt, 0), g.domain_id, d.name, d.status
						FROM user_groups g LEFT JOIN domains d ON g.domain_id=d.id
						LEFT JOIN (SELECT group_id, COUNT(*) AS cnt FROM group_members GROUP BY group_id) c ON g.id=c.group_id
						WHERE g.domain_id=$1 AND g.status=$2
						ORDER BY g.name`
	rows, err := db.Query(query, domainId, STATUS_ACTIVE)
	if err != nil {
		return nil
	}

	defer rows.Close()

	var groups []*model.Group2

	for {
		group := readGroupRow(rows)
		if group == nil {
			break
		}
		groups = append(groups, group)
	}

	return groups
}

// Select group with the name / id
func SelectGroup(domainId int, groupName string, groupId int) *model.Group2 {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	if groupId > 0 {
		query = `SELECT g.id, g.name, COALESCE(c.cnt, 0), g.domain_id, d.name, d.status
						FROM user_groups g LEFT JOIN domains d ON g.domain_id=d.id
						LEFT JOIN (SELECT group_id, COUNT(*) AS cnt FROM group_members GROUP BY group_id) c ON g.id=c.group_id
						WHERE g.domain_id=$1 AND g.status=$2 AND g.id=$3`
		rows, err = db.Query(query, domainId, STATUS_ACTIVE, groupId)
	} else {
		query = `SELECT g.id, g.name, COALESCE(c.cnt, 0), g.domain_id, d.name, d.status
						FROM user_groups g LEFT JOIN domains d ON g.domain_id=d.id
						LEFT JOIN (SELECT group_id, COUNT(*) AS cnt FROM group_members GROUP BY group_id) c ON g.id=c.group_id
						WHERE g.domain_id=$1 AND g.status=$2 AND g.name=$3`
		rows, err = db.Query(query, domainId, STATUS_ACTIVE, groupName)
	}

	fmt.Printf("SelectGroup: domain=%d name=%s id=%d err=%v\n", domainId, groupName, groupId, err)

	if err != nil {
		return nil
	}

	group := readGroupRow(rows)

	rows.Close()

	return group
}

// Update the the record with the id
func UpdateGroup(domainId int, groupName string, groupId int, group *model.Group) *model.Group2 {
	db := setupDB()

	// Compose SQL query
	var query string
	var err error
	var rowsAffected int64
	var result sql.Result
	var params string

	if group.Name != "" {
		params += "name='" + group.Name + "', "
	}
	if params == "" {
		// Nothing to update
		return nil
	}
	params = params[:len(params)-2]

	if groupId > 0 {
		query = fmt.Sprintf("UPDATE user_groups SET %s WHERE domain_id=$1 AND id=$2 AND status=$3", params)
		result, err = db.Exec(query, domainId, groupId, STATUS_ACTIVE)
	} else {
		query = fmt.Sprintf("UPDATE user_groups SET %s WHERE domain_id=$1 AND name=$2 AND status=$3", params)
		result, err = db.Exec(query, domainId, groupName, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("Update Group %s %d domain %d - %v\n",
			groupName, groupId, domainId, err)
		return nil
	}
	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("Update Group %s %d domain %d - %d rows affected\n",
		groupName, groupId, domainId, rowsAffected)

	// if rowsAffected != 1 {
	// 	// Unexpected
	// }

	// Select the updated record and return
	return SelectGroup(domainId, groupName, groupId)
}

// Return all users of the specified group
func GetGroupUsers(domainId int, groupName string, groupId int) []*model.User2 {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	var resp []*model.User2

	if groupId > 0 {
		query = `SELECT u.id, u.name, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
						FROM user_groups g, users u, group_members members
						LEFT JOIN domains d ON g.domain_id=d.id
						WHERE u.id=members.user_id
							AND g.domain_id=$1
							AND g.status=$2
							AND members.group_id=g.id
							AND g.id=$3`
		rows, err = db.Query(query, domainId, STATUS_ACTIVE, groupId)
	} else {
		query = `SELECT u.id, u.name, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
						FROM user_groups g, users u, group_members members
						LEFT JOIN domains d ON g.domain_id=d.id
						WHERE u.id=members.user_id
							AND g.domain_id=$1
							AND g.status=$2
							AND members.group_id=g.id
							AND g.name=$3`
		rows, err = db.Query(query, domainId, STATUS_ACTIVE, groupName)
	}
	if err != nil {
		return nil
	}
	for {
		user := readUserRow(rows, true)
		if user == nil {
			break
		}
		resp = append(resp, user)
	}
	rows.Close()
	return resp
}

// Delete the record with the id
func DeleteGroup(domainId int, groupName string, groupId int) (*model.Group2, int) {
	db := setupDB()

	deleted_group := SelectGroup(domainId, groupName, groupId)
	if deleted_group == nil {
		// specified group not found
		return nil, 0
	}

	if groupId == 0 {
		groupId = deleted_group.ID
	}

	// if groupName == "" {
	// 	groupName = deleted_group.Name
	// }

	var err error
	_, err = db.Exec("UPDATE user_groups SET status=$1 WHERE id=$2 AND domain_id=$3",
		STATUS_DELETED, groupId, domainId)
	if err != nil {
		return nil, 0
	}

	//
	// Delete all group memebers
	//
	rowsAffected := RemoveAllMembers(groupId)

	return deleted_group, rowsAffected
}

func readGroupRow(rows *sql.Rows) *model.Group2 {
	var id int
	var name sql.NullString
	var domainId int
	var cnt int
	var dname sql.NullString
	var dstatus sql.NullString

	if !rows.Next() {
		return nil
	}

	err := rows.Scan(&id, &name, &cnt, &domainId, &dname, &dstatus)
	if err != nil {
		return nil
	}

	domain := model.Domain{ID: domainId, Name: dname.String, Status: dstatus.String}
	group := model.Group2{ID: id, Name: name.String, Count: cnt, Domain: domain}
	return &group
}
