package db

import (
	"database/sql"
	"fmt"

	model "github.com/saroopmathur/rest-api/models"
)

func addGroupMember(domainId int, groupName string, groupId int, username string) error {
	db := setupDB()

	var query string
	var err error
	var lastInsertID int
	if groupId > 0 {
		query = `INSERT INTO group_members (group_id, user_id)
				SELECT g.id, u.id from user_groups g, users u
					WHERE u.name=$1 AND g.id=$2
						AND g.domain_id=$3 AND u.domain_id=$3
						AND u.status=$4 AND g.status=$4 returning id`
		err = db.QueryRow(query, username, groupId, domainId, STATUS_ACTIVE).Scan(&lastInsertID)
	} else {
		query = `INSERT INTO group_members (group_id, user_id)
				SELECT g.id, u.id from user_groups g, users u
					WHERE u.name=$1 AND g.name=$2
						AND g.domain_id=$3 AND u.domain_id=$3
						AND u.status=$4 AND g.status=$4 returning id`
		err = db.QueryRow(query, username, groupName, domainId, STATUS_ACTIVE).Scan(&lastInsertID)
	}
	fmt.Printf("addGroupMember [%s %d %d] %s %v\n", groupName, groupId, domainId, username, err)
	if err != nil {
		return err
	}

	return nil
}

func AddGroupMembers(domainId int, groupName string, groupId int, users []string) int {
	var addCount int
	var err error
	for _, userName := range users {
		err = addGroupMember(domainId, groupName, groupId, userName)
		if err == nil {
			addCount++
		}
	}
	return addCount
}

func removeGroupMemberByName(domainId int, groupName string, groupId int, userName string) (int, error) {
	var result sql.Result
	var err error
	var query string

	db := setupDB()
	if groupId > 0 {
		query = `DELETE FROM group_members WHERE group_id=$1
				AND user_id=(SELECT id from users WHERE name=$2 AND domain_id=$3)`
		result, err = db.Exec(query, groupId, userName, domainId)
	} else {
		query = `DELETE FROM group_members WHERE
				group_id=(SELECT id from user_groups WHERE name=$1 AND domain_id=$2)
				AND user_id=(SELECT id from users WHERE name=$3 AND domain_id=$4)`
		result, err = db.Exec(query, groupName, domainId, userName, domainId)
	}

	if err != nil {
		fmt.Printf("Remove user %s from group [%s %d] Domain:%d - %v\n",
			userName, groupName, groupId, domainId, err)
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Removed user %s from group [%s %d] Domain:%d - %d rows affected\n",
		userName, groupName, groupId, domainId, rowsAffected)
	return int(rowsAffected), nil
}

func RemoveGroupMembers(domainId int, groupName string, groupId int, users []string) int {
	var rowsAffected int
	for _, userName := range users {
		n, err := removeGroupMemberByName(domainId, groupName, groupId, userName)
		if err != nil {
			continue
		}
		rowsAffected += n
	}
	return rowsAffected
}

func SelectGroupMembers(domainId int, groupName string, groupId int) []*model.User2 {
	var rows *sql.Rows
	var err error
	var query string

	db := setupDB()
	if groupId == 0 {
		// Select by group name
		query = `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status, g.id, g.name
				FROM user_groups g, group_members mem, users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE g.name=$1
					AND mem.group_id=g.id 
					AND mem.user_id=u.id
					AND u.domain_id=$2 AND u.status=$3`
		rows, err = db.Query(query, groupName, domainId, STATUS_ACTIVE)
	} else {
		// Select by group id
		query = `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status, g.id, g.name
				FROM user_groups g, group_members mem, users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE g.id=$1
					AND mem.group_id=g.id 
					AND mem.user_id=u.id
					AND u.domain_id=$2 AND u.status=$3`
		rows, err = db.Query(query, groupId, domainId, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("SelectGroupMembers: [%s %d %d] %v\n", groupName, groupId, domainId, err)
		return nil
	}
	defer rows.Close()

	var users []*model.User2
	for {
		user := readUserRow(rows, true)
		if user == nil {
			break
		}
		users = append(users, user)
	}
	fmt.Printf("SelectGroupMembers: [%s %d %d] Read %d members\n", groupName, groupId, domainId, len(users))
	return users
}

func GetUserGroups(domainId int, userName string, userId int) []*model.Group2 {
	var rows *sql.Rows
	var err error
	var query string

	db := setupDB()

	if userId == 0 {
		// Lookup by userName
		query = `SELECT g.id, g.name, 0, g.domain_id, d.name, d.status
				FROM user_groups g, group_members member, users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE u.name=$1
					AND member.user_id=u.id
					AND g.id=member.group_id
					AND g.domain_id=$2
					AND g.status=$3
				ORDER BY g.name`
		rows, err = db.Query(query, userName, domainId, STATUS_ACTIVE)
	} else {
		// Lookup by userId
		query = `SELECT g.id, g.name, 0, g.domain_id, d.name, d.status
				FROM group_members member, user_groups g LEFT JOIN domains d ON g.domain_id=d.id
				WHERE member.user_id=$1
					AND g.id=member.group_id
					AND g.domain_id=$2
					AND g.status=$3
				ORDER BY g.name`
		rows, err = db.Query(query, userId, domainId, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("GetUserGroups: [%s %d %d] %v\n", userName, userId, domainId, err)
		return nil
	}

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

// Remove all members of specified group groupId
// Must only be called after verifying that the admin has permissions to this group
func RemoveAllMembers(groupId int) int {
	db := setupDB()

	result, err := db.Exec("DELETE FROM group_members WHERE group_id=$1", groupId)
	if err != nil {
		fmt.Printf("Deleted all members of group %d err=%v\n", groupId, err)
		return 0
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Deleted all members of group %d - %d rows affected\n", groupId, rowsAffected)
	return int(rowsAffected)
}
