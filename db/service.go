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
func InsertService(domainId int, service *model.Service) (*model.Service2, error) {
	db := setupDB()

	// Delete if any record with the same name
	query := `DELETE FROM services WHERE domain_id=$1 AND name=$2`
	db.Exec(query, domainId, service.Name)

	var lastInsertID int
	query = `INSERT INTO services (domain_id, name, password, wg_key, status)
						VALUES ($1, $2, $3, $4, $5) returning id`
	err := db.QueryRow(query, domainId, service.Name, service.Password, service.WGKey, STATUS_ACTIVE).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	// Create an entry in apps table
	query = `INSERT INTO apps (name, service_id, status)
						VALUES ($1, $2, $3)`
	db.Exec(query, service.Name, lastInsertID, STATUS_ACTIVE)

	// Select the inserted record and return
	return SelectService(domainId, "", lastInsertID), err
}

// Select returns the whole database
func SelectServices(domainId int) []*model.Service2 {
	db := setupDB()

	query := `SELECT s.id, s.name, s.password, COALESCE(s.icon, 0), s.wg_key, s.local_ip, s.public_ip, s.virtual_ip, d.id, d.name, d.status
			FROM services s LEFT JOIN domains d ON s.domain_id=d.id
			WHERE s.domain_id=$1 AND s.status=$2 ORDER BY s.name`
	rows, err := db.Query(query, domainId, STATUS_ACTIVE)
	if err != nil {
		return nil
	}

	var services []*model.Service2

	for {
		service := readServiceRow(rows)
		if service == nil {
			break
		}
		services = append(services, service)
	}
	return services
}

// Select the service with either name or id
func SelectService(domainId int, serviceName string, serviceId int) *model.Service2 {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	if serviceId == 0 {
		query = `SELECT s.id, s.name, s.password, COALESCE(s.icon, 0), s.wg_key, s.local_ip, s.public_ip, s.virtual_ip, d.id, d.name, d.status
				FROM services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND s.name=$2 AND s.status=$3`
		rows, err = db.Query(query, domainId, serviceName, STATUS_ACTIVE)
	} else {
		query = `SELECT s.id, s.name, s.password, COALESCE(s.icon, 0), s.wg_key, s.local_ip, s.public_ip, s.virtual_ip, d.id, d.name, d.status
				FROM services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND s.id=$2 AND s.status=$3`
		rows, err = db.Query(query, domainId, serviceId, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("%s: domain=%d [%s %d] %v\n", query, domainId, serviceName, serviceId, err)
		return nil
	}

	return readServiceRow(rows)
}

// Update the the record with the id
func UpdateService(domainId int, serviceName string, serviceId int, service *model.Service) *model.Service2 {
	db := setupDB()

	// Compose SQL query
	var params string
	if service.Name != "" {
		params += "name='" + service.Name + "', "
	}
	if service.Password != "" {
		params += "password='" + service.Password + "', "
	}
	if service.WGKey != "" {
		params += "wg_key='" + service.WGKey + "', "
	}
	if service.LocalIP != "" {
		params += "local_ip='" + service.LocalIP + "', "
	}
	if service.PublicIP != "" {
		params += "public_ip='" + service.PublicIP + "', "
	}
	if service.VirtualIP != "" {
		params += "virtual_ip='" + service.VirtualIP + "', "
	}
	if params == "" {
		// Nothing to update
		return nil
	}
	params = params[:len(params)-2]

	var err error
	var query string
	if serviceId > 0 {
		query = fmt.Sprintf("UPDATE services SET %s WHERE id=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, serviceId, domainId, STATUS_ACTIVE)
	} else {
		query = fmt.Sprintf("UPDATE services SET %s WHERE name=$1 AND domain_id=$2 AND status=$3", params)
		_, err = db.Exec(query, serviceName, domainId, STATUS_ACTIVE)
	}
	// fmt.Printf("%s domainId=%d serviceId=%d serviceName=%s err=%v\n", query, domainId, serviceId, serviceName, err)

	if err != nil {
		fmt.Printf("%s: [%s %d] %v\n", query, serviceName, serviceId, err)
		return nil
	}

	// Select the updated record and return
	return SelectService(domainId, serviceName, serviceId)
}

// Delete the record with the id
func DeleteService(domainId int, serviceName string, serviceId int) *model.Service2 {
	db := setupDB()

	deleted_service := SelectService(domainId, serviceName, serviceId)
	if deleted_service == nil {
		// Invalid serviceName or serviceId
		return nil
	}

	// Delete all the records in apps with the service_id
	query := "UPDATE apps set status=$1 WHERE service_id=$2"
	_, err := db.Exec(query, STATUS_DELETED, deleted_service.ID)
	if err != nil {
		return nil
	}

	query = "UPDATE services SET status=$1 WHERE id=$2 AND domain_id=$3"
	_, err = db.Exec(query, STATUS_DELETED, deleted_service.ID, domainId)

	fmt.Printf("%s domainId=%d serviceId=%d serviceName=%s err=%v\n", query, domainId, serviceId, serviceName, err)
	if err != nil {
		fmt.Printf("%s: [%s %d] %v\n", query, serviceName, serviceId, err)
		return nil
	}

	return deleted_service
}

func GenerateAndSaveServiceToken(s *model.Service2) {
	db := setupDB()

	token := fmt.Sprintf("S%d", rand.Int63())

	query := `INSERT INTO sessions (uid, session_id, domain_id, role, start_time, status)
				     VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Query(query, s.ID, token, s.Domain.ID, ROLE_SERVICE, time.Now(), STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("%d: %s\n", s.ID, token)
		return
	}
	s.SessionID = token
}

func GetServiceByToken(token string) *model.Service2 {
	if !strings.HasPrefix(token, "S") {
		return nil
	}
	fmt.Printf("GetServiceByToken: %s\n", token)
	db := setupDB()

	query := `SELECT s.id, s.name, s.password, s.icon, s.wg_key, s.local_ip, s.public_ip, s.virtual_ip, d.id, d.name, d.status
				FROM sessions sess, services s, domains d
				WHERE sess.session_id=$1
					AND s.domain_id=d.id
					AND s.id=sess.uid
					AND sess.role=$2
					AND sess.status=$3`

	//fmt.Printf("%s: [%s]\n", query, token)
	rows, err := db.Query(query, token, ROLE_SERVICE, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("%s\n", token)
		return nil
	}
	defer rows.Close()

	return readServiceRow(rows)
}

func GetServiceByName(servicename string) *model.Service2 {
	fmt.Printf("GetServiceByName: %s\n", servicename)
	parts := strings.Split(servicename, "@")
	if len(parts) != 2 {
		// missing domain name
		return nil
	}
	name := parts[0]
	domain := parts[1]

	db := setupDB()

	query := `SELECT s.id, s.name, s.password, s.icon, s.wg_key, s.local_ip, s.public_ip, s.virtual_ip, d.id, d.name, d.status
				FROM services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.name=$1 AND d.name=$2 AND s.status=$3`
	//fmt.Printf("%s: [%s@%s]\n", query, name, domain)
	rows, err := db.Query(query, name, domain, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("[%s@%s] DB Query Failed\n", name, domain)
		return nil
	}
	defer rows.Close()

	return readServiceRow(rows)
}

func readServiceRow(rows *sql.Rows) *model.Service2 {
	var serviceId int
	var domainId int
	var name sql.NullString
	var password sql.NullString
	var icon int
	var dname sql.NullString
	var dstatus sql.NullString
	var wgKey sql.NullString
	var localIp sql.NullString
	var publicIp sql.NullString
	var virtualIp sql.NullString

	if !rows.Next() {
		//fmt.Printf("ReadService Next is false\n")
		return nil
	}

	err := rows.Scan(&serviceId, &name, &password, &icon, &wgKey, &localIp, &publicIp, &virtualIp, &domainId, &dname, &dstatus)
	if err != nil {
		fmt.Printf("ReadService Scan: %v\n", err)
		return nil
	}

	domain := model.Domain{ID: domainId, Name: dname.String, Status: dstatus.String}
	service := model.Service2{
		ID:        serviceId,
		Name:      name.String,
		Password:  password.String,
		Icon:      icon,
		Domain:    domain,
		WGKey:     wgKey.String,
		LocalIP:   localIp.String,
		PublicIP:  publicIp.String,
		VirtualIP: virtualIp.String,
	}
	//fmt.Printf("ReadService: %s@%s\n", service.Name, domain.Name)
	return &service
}
