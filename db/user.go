package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	model "github.com/saroopmathur/rest-api/models"
)

// Insert allows populating database
func InsertUser(domainId int, user *model.User) (*model.User2, error) {
	db := setupDB()

	// Delete if any record with the same name
	query := `DELETE FROM users WHERE domain_id=$1 AND name=$2`
	db.Exec(query, domainId, user.Name)

	var lastInsertID int
	query = `INSERT INTO users (domain_id, name, password, wg_key, status)
						VALUES ($1, $2, $3, $4, $5) returning id`
	err := db.QueryRow(query, domainId, user.Name, user.Password, user.WGKey, STATUS_ACTIVE).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Created New User %s in domain %d\n", user.Name, domainId)

	// Select the inserted record and return
	return SelectUser(domainId, "", lastInsertID), err
}

// Select returns the whole database
func SelectUsers(domainId int, offset int, limit int, search string) []*model.User2 {
	db := setupDB()

	query := `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id AS did, d.name AS dname, d.status
				FROM users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE u.domain_id=$1 AND u.status=$2 AND u.name LIKE $3 ORDER BY u.name OFFSET $4 LIMIT $5`
	rows, err := db.Query(query, domainId, STATUS_ACTIVE, search+"%", offset, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []*model.User2

	for {
		user := readUserRow(rows, false)
		if user == nil {
			break
		}
		users = append(users, user)
	}
	return users
}

// Select the user with either name or id
func SelectUser(domainId int, userName string, userId int) *model.User2 {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	var user *model.User2

	if userId == 0 {
		query = `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
				FROM users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE u.domain_id=$1 AND u.name=$2 AND u.status=$3`
		rows, err = db.Query(query, domainId, userName, STATUS_ACTIVE)
	} else {
		query = `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
				FROM users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE u.domain_id=$1 AND u.id=$2 AND u.status=$3`
		rows, err = db.Query(query, domainId, userId, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("%s: domain=%d [%s %d] %v\n", query, domainId, userName, userId, err)
		return nil
	}
	user = readUserRow(rows, false)
	if user != nil {
		//fmt.Printf("Select User %s %d in domain %d - %v\n", userName, userId, domainId, *user)
	} else {
		fmt.Printf("Select User %s %d in domain %d - NOT FOUND\n", userName, userId, domainId)
	}
	rows.Close()
	return user
}

// Update the the record with the id
func UpdateUser(domainId int, userName string, userId int, user *model.User) *model.User2 {
	db := setupDB()

	// Compose SQL query
	var params string
	if user.Name != "" {
		params += "name='" + user.Name + "', "
	}
	if user.Password != "" {
		params += "password='" + user.Password + "', "
	}
	if user.WGKey != "" {
		params += "wg_key='" + user.WGKey + "', "
	}
	if user.LocalIP != "" {
		params += "local_ip='" + user.LocalIP + "', "
	}
	if user.PublicIP != "" {
		params += "public_ip='" + user.PublicIP + "', "
	}
	if user.VirtualIP != "" {
		params += "virtual_ip='" + user.VirtualIP + "', "
	}
	if params == "" {
		// Nothing to update
		return nil
	}
	params = params[:len(params)-2]

	var err error
	var query string
	if userId > 0 {
		query = fmt.Sprintf("UPDATE users SET %s WHERE id=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, userId, domainId, STATUS_ACTIVE)
	} else {
		query = fmt.Sprintf("UPDATE users SET %s WHERE name=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, userName, domainId, STATUS_ACTIVE)
	}
	//fmt.Printf("%s domainId=%d userId=%d userName=%s err=%v\n", query, domainId, userId, userName, err)

	if err != nil {
		fmt.Printf("%s: [%s %d] %v\n", query, userName, userId, err)
		return nil
	}

	// Select the updated record and return
	return SelectUser(domainId, userName, userId)
}

// Delete the record with the id
func DeleteUser(domainId int, userName string, userId int) *model.User2 {
	db := setupDB()

	deleted_user := SelectUser(domainId, userName, userId)
	if deleted_user == nil {
		// Invalid userName or userId
		fmt.Printf("DeleteUser: [%s %d] domain %d - Invalid User\n", userName, userId, domainId)
		return nil
	}

	var err error
	var query string
	if userId > 0 {
		query = "UPDATE users SET status=$1 WHERE id=$2 AND domain_id=$3"
		_, err = db.Exec(query, STATUS_DELETED, userId, domainId)
	} else {
		query = "UPDATE users SET status=$1 WHERE name=$2 AND domain_id=$3"
		_, err = db.Exec(query, STATUS_DELETED, userName, domainId)
	}
	if err != nil {
		fmt.Printf("DeleteUser: [%s %d] domain %d - %v\n", userName, userId, domainId, err)
		return nil
	}

	return deleted_user
}

func GenerateAndSaveUserToken(u *model.User2) {
	db := setupDB()

	token := fmt.Sprintf("U%d", rand.Int63())

	query := `INSERT INTO sessions (uid, session_id, domain_id, role, start_time, status)
				     VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Query(query, u.ID, token, u.Domain.ID, ROLE_USER, time.Now(), STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("GenerateToken: %s %v\n", token, err)
		return
	}
	u.SessionID = token
}

func GetUserByToken(token string) *model.User2 {
	if !strings.HasPrefix(token, "U") {
		return nil
	}
	//fmt.Printf("GetUserByToken: %s\n", token)
	db := setupDB()

	query := `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
				FROM sessions sess, users u, domains d
				WHERE sess.session_id=$1
					AND sess.uid=u.id
					AND u.domain_id=d.id
					AND sess.status=$2`

	//fmt.Printf("%s: [%s]\n", query, token)
	rows, err := db.Query(query, token, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("%s\n", token)
		return nil
	}
	defer rows.Close()

	return readUserRow(rows, false)
}

func GetUserByName(username string) *model.User2 {
	//fmt.Printf("GetUserByName: %s\n", username)
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		// missing domain name
		return nil
	}
	name := parts[0]
	domain := parts[1]

	db := setupDB()

	query := `SELECT u.id, u.name, u.password, u.wg_key, u.local_ip, u.public_ip, u.virtual_ip, d.id, d.name, d.status
				FROM users u LEFT JOIN domains d ON u.domain_id=d.id
				WHERE u.name=$1 AND d.name=$2 AND u.status=$3`
	//fmt.Printf("%s: [%s@%s]\n", query, name, domain)
	rows, err := db.Query(query, name, domain, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("[%s@%s] DB Query Failed\n", name, domain)
		return nil
	}
	defer rows.Close()

	return readUserRow(rows, false)
}

func readUserRow(rows *sql.Rows, readGroup bool) *model.User2 {
	var userId int
	var name sql.NullString
	var pass sql.NullString
	var domainId int
	var dname sql.NullString
	var dstatus sql.NullString
	var groupId sql.NullInt32
	var gname sql.NullString
	var wgKey sql.NullString
	var localIp sql.NullString
	var publicIp sql.NullString
	var virtualIp sql.NullString

	if !rows.Next() {
		return nil
	}

	var err error
	if readGroup {
		err = rows.Scan(&userId, &name, &pass, &wgKey, &localIp, &publicIp, &virtualIp, &domainId, &dname, &dstatus, &groupId, &gname)
	} else {
		err = rows.Scan(&userId, &name, &pass, &wgKey, &localIp, &publicIp, &virtualIp, &domainId, &dname, &dstatus)
	}
	if err != nil {
		fmt.Printf("ReadUser Scan: %v\n", err)
		return nil
	}

	domain := model.Domain{ID: domainId, Name: dname.String, Status: dstatus.String}
	group := model.Group{ID: int(groupId.Int32), Name: gname.String}
	user := model.User2{
		ID:        userId,
		Name:      name.String,
		Password:  pass.String,
		Domain:    domain,
		Group:     group,
		WGKey:     wgKey.String,
		LocalIP:   localIp.String,
		PublicIP:  publicIp.String,
		VirtualIP: virtualIp.String,
	}
	//fmt.Printf("ReadUser: %s\n", user.Name)
	return &user
}
