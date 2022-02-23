package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	model "github.com/saroopmathur/rest-api/models"
)

// Create allows populating database
func InsertAdmin(domainId int, name string, password string) (*model.Admin2, error) {
	db := setupDB()

	// Delete if any record with the same name
	query := `DELETE FROM admins WHERE name=$1`
	db.Exec(query, name)

	var lastInsertID int
	query = "INSERT INTO admins (name, domain_id, password, status) VALUES($1, $2, $3, $4) returning id"
	err := db.QueryRow(query, name, domainId, password, STATUS_ACTIVE).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	// Select the inserted record and return
	return SelectAdmin(domainId, "", lastInsertID), nil
}

// Select returns the whole database
func SelectAdmins(domainId int) []*model.Admin2 {
	db := setupDB()

	query := `SELECT a.id, a.name, a.domain_id, a.password, d.name, d.status
						FROM admins a LEFT JOIN domains d ON a.domain_id=d.id
						WHERE a.domain_id=$1 AND a.status=$2 ORDER BY a.name`
	rows, err := db.Query(query, domainId, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("SelectAdmins: %v\n", err)
		return nil
	}
	defer rows.Close()

	var admins []*model.Admin2

	for {
		admin := readAdminRow(rows, false)
		if admin == nil {
			break
		}
		admins = append(admins, admin)
	}
	return admins
}

// Select the record with the id
func SelectAdmin(domain_id int, adminName string, adminId int) *model.Admin2 {
	db := setupDB()

	var rows *sql.Rows
	var err error
	if adminId == 0 {
		query := `SELECT a.id, a.name, a.domain_id, a.password, d.name, d.status
						FROM admins a LEFT JOIN domains d ON a.domain_id=d.id
						WHERE a.name=$1 AND a.domain_id=$2 AND a.status=$3`
		rows, err = db.Query(query, adminName, domain_id, STATUS_ACTIVE)
	} else {
		query := `SELECT a.id, a.name, a.domain_id, a.password, d.name, d.status
						FROM admins a LEFT JOIN domains d ON a.domain_id=d.id
						WHERE a.id=$1 AND a.domain_id=$2 AND a.status=$3`
		rows, err = db.Query(query, adminId, domain_id, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("SelectAdmin: %v\n", err)
		return nil
	}

	admin := readAdminRow(rows, false)
	rows.Close()
	return admin
}

// Update the the record with the id
func UpdateAdmin(domainId int, adminName string, adminId int, admin *model.Admin) *model.Admin2 {
	db := setupDB()

	name := admin.Name
	pass := admin.Password

	// Compose SQL query
	var params string
	if name != "" {
		params += "name='" + name + "', "
	}
	if pass != "" {
		params += "password='" + pass + "', "
	}
	if params == "" {
		// Nothing to update
		return nil
	}
	params = params[:len(params)-2]

	var err error
	if adminId > 0 {
		query := fmt.Sprintf("UPDATE admins SET %s WHERE id=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, adminId, domainId, STATUS_ACTIVE)
	} else {
		query := fmt.Sprintf("UPDATE admins SET %s WHERE name=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, adminName, domainId, STATUS_ACTIVE)
	}

	if err != nil {
		return nil
	}

	// Select the updated record and return
	return SelectAdmin(domainId, adminName, adminId)
}

// Delete the record with the adminId or adminName
// domain_id is the domain for account that is atempting this operation
func DeleteAdmin(domainId int, adminName string, adminId int) error {
	db := setupDB()
	var err error
	if adminId != 0 {
		// delete by ID
		_, err = db.Exec("UPDATE admins SET status=$3 WHERE id=$1 AND domain_id=$2",
			adminId, domainId, STATUS_DELETED)
	} else {
		// delete by name
		_, err = db.Exec("UPDATE admins SET status=$3 WHERE name=$1 AND domain_id=$2",
			adminName, domainId, STATUS_DELETED)
	}

	return err
}

func GetAdminByName(username string) *model.Admin2 {
	//log.Printf("GetAdminByName: %s\n", username)
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		// missing domain name
		return nil
	}
	name := parts[0]
	domain := parts[1]

	db := setupDB()

	query := `SELECT a.id, a.name, a.domain_id, a.password, d.name, d.status
				FROM admins a LEFT JOIN domains d ON a.domain_id=d.id
				WHERE a.name=$1 AND d.name=$2 AND a.status=$3`
	rows, err := db.Query(query, name, domain, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("GetAdminByName: %s Failed\n", username)
		return nil
	}
	admin := readAdminRow(rows, false)
	rows.Close()

	if admin == nil {
		fmt.Printf("GetAdminByName: [%s@%s] \"%s\" Not Found\n", name, domain, query)
	} else {
		fmt.Printf("GetAdminByName: %s %+v\n", username, *admin)
	}
	return admin
}

func readAdminRow(rows *sql.Rows, readRole bool) *model.Admin2 {
	var admin *model.Admin2
	if rows.Next() {
		var id int
		var name string
		var domainId int
		var dname sql.NullString
		var dstatus sql.NullString
		var pass sql.NullString
		var role sql.NullString
		var err error

		if readRole {
			err = rows.Scan(&id, &name, &domainId, &pass, &role, &dname, &dstatus)
		} else {
			err = rows.Scan(&id, &name, &domainId, &pass, &dname, &dstatus)
			role.String = ROLE_ADMIN
		}
		if err != nil {
			return nil
		}

		domain := model.Domain{ID: domainId, Name: dname.String, Status: dstatus.String}
		admin = &model.Admin2{ID: id, Name: name, Domain: domain, Role: role.String, Password: pass.String}
		//log.Printf("AdminfromDB: %v\n", admin)
	} else {
		log.Printf("AdminfromDB: Not Found\n")
		return nil
	}

	return admin
}

func GenerateAndSaveAdminToken(a *model.Admin2) {
	db := setupDB()

	token := fmt.Sprintf("A%d", rand.Int63())

	query := `INSERT INTO sessions (uid, session_id, domain_id, role, start_time, status)
				     VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Query(query, a.ID, token, a.Domain.ID, a.Role, time.Now(), STATUS_ACTIVE)
	if err != nil {
		log.Printf("SaveAdminToken Failed: %s %d %v\n", a.Name, a.ID, err)
		return
	}
	a.SessionID = token
	//fmt.Printf("SaveAdminToken: role=%s token=%s [%s %d]\n", a.Role, token, a.Name, a.ID)
}

func GetAdminByToken(token string) *model.Admin2 {
	if !strings.HasPrefix(token, "A") {
		return nil
	}
	db := setupDB()

	query := `SELECT a.id, a.name, sess.domain_id, a.password, sess.role, d.name, d.status
				FROM sessions sess, admins a, domains d
				WHERE sess.session_id=$1
					AND sess.uid=a.id
					AND sess.domain_id=d.id
					AND sess.status=$2
					AND a.status=$2`
	//log.Printf("%s: [%s]\n", query, token)
	rows, err := db.Query(query, token, STATUS_ACTIVE)
	if err != nil {
		log.Printf("GetAdminByToken: %s Failed\n", token)
		return nil
	}

	admin := readAdminRow(rows, true)
	rows.Close()

	//fmt.Printf("ReadAdminToken: role=%s token=%s [%s %d]\n", admin.Role, token, admin.Name, admin.ID)
	return admin
}
