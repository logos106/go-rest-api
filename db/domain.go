package db

import (
	"database/sql"
	"fmt"
	"log"

	model "github.com/saroopmathur/rest-api/models"
)

// Insert allows populating database
func InsertDomain(domain *model.DomainReq) (*model.Domain, error) {
	db := setupDB()

	var lastInsertID int
	query := `INSERT INTO domains (name, status) VALUES($1, $2) returning id`
	err := db.QueryRow(query, domain.Name, STATUS_ACTIVE).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	// Select the inserted record and return
	return SelectDomain(lastInsertID, domain.Name), nil
}

// Select returns the whole database
func SelectDomains() []*model.Domain {
	db := setupDB()

	query := "SELECT id, name FROM domains WHERE status=$1 ORDER BY name"
	rows, err := db.Query(query, STATUS_ACTIVE)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var domains []*model.Domain

	for rows.Next() {
		var id int
		var name sql.NullString

		err = rows.Scan(&id, &name)
		if err != nil {
			return nil
		}
		domains = append(domains, &model.Domain{ID: id, Name: name.String})
	}

	return domains
}

// Select the record with the id
func SelectDomain(domainId int, domainName string) *model.Domain {
	db := setupDB()
	var rows *sql.Rows
	var err error

	if domainId > 0 {
		rows, err = db.Query("SELECT id, name FROM domains WHERE id=$1 AND status=$2", domainId, STATUS_ACTIVE)
	} else {
		rows, err = db.Query("SELECT id, name FROM domains WHERE name=$1 AND status=$2", domainName, STATUS_ACTIVE)
	}
	if err != nil {
		return nil
	}
	defer rows.Close()

	var domain *model.Domain

	var id int
	var name sql.NullString

	if rows.Next() {
		err = rows.Scan(&id, &name)
		if err == nil {
			domain = &model.Domain{ID: id, Name: name.String}
		}
	}

	return domain
}

// Update the the record with the id
func UpdateDomain(domainId int, domainName string, domain *model.DomainReq) *model.Domain {
	db := setupDB()

	name := domain.Name

	if name == "" {
		// Nothing to do
		return SelectDomain(domainId, name)
	}

	// Compose SQL query
	var err error
	var query string
	if domainId > 0 {
		query = "UPDATE domains SET name=$1 WHERE id=$2 AND status=$3"
		_, err = db.Exec(query, name, domainId, STATUS_ACTIVE)
	} else {
		query = "UPDATE domains SET name=$1 WHERE name=$2 AND status=$3"
		_, err = db.Exec(query, name, domainName, STATUS_ACTIVE)
	}
	if err != nil {
		return nil
	}

	// Select the updated record and return
	return SelectDomain(domainId, name)
}

// Delete the record with the id
func DeleteDomain(domainId int, domainName string) *model.Domain {
	db := setupDB()

	deleted_domain := SelectDomain(domainId, domainName)
	if deleted_domain == nil {
		// Unknown domain name or Id
		return nil
	}

	// Try to delete one to check if any foreign key exists
	query := "DELETE FROM domains WHERE id=$1 OR name=$2"
	_, err := db.Exec(query, domainId, domainName)
	if err != nil {
		return nil
	}

	// If deleted, restore it and update status
	query = "INSERT INTO domains (id, name, status) VALUES ($1, $2, $3)"
	_, err = db.Exec(query, deleted_domain.ID, deleted_domain.Name, STATUS_DELETED)
	if err != nil {
		fmt.Printf("%s: [%s %d] %v\n", query, domainName, domainId, err)
		return nil
	}

	// _, err = db.Exec("UPDATE domains SET status=$1 WHERE id=$2", STATUS_DELETED, deleted_domain.ID)
	// if err != nil {
	// 	return nil
	// }

	return deleted_domain
}

func ChangeDomain(sessionId string, newDomainId int, newDomainName string) error {
	db := setupDB()
	var err error

	if newDomainId > 0 {
		// Updated based on domain id
		query := `UPDATE sessions SET domain_id=$1 WHERE session_id=$2 AND status=$3`
		_, err = db.Query(query, newDomainId, sessionId, STATUS_ACTIVE)
	} else {
		// Updated based on domain name
		query := `UPDATE sessions SET domain_id=(SELECT id FROM domains WHERE name=$1) WHERE session_id=$2 AND status=$3`
		_, err = db.Query(query, newDomainName, sessionId, STATUS_ACTIVE)
	}
	if err != nil {
		log.Printf("ChangeDomain: %s [%s %d] %v\n", sessionId, newDomainName, newDomainId, err)
	} else {
		log.Printf("ChangeDomain: %s [%s %d] SUCCESS\n", sessionId, newDomainName, newDomainId)
	}
	return nil
}

// Set admin to any domain, other than powerdomain
// If there are no other domains, then set to powerdomain
func SetAnyDomain(dom *model.Domain) {
	var selected *model.Domain

	domains := SelectDomains()
	for _, domain := range domains {
		if domain.ID == POWERDOMAIN {
			selected = domain
		} else {
			selected = domain
			break
		}
	}

	if selected != nil {
		*dom = *selected
	}
}
