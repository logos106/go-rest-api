package db

import (
	"database/sql"
	"fmt"
	"log"

	model "github.com/saroopmathur/rest-api/models"
)

// Insert allows populating database
func InsertApp(domainId int, app *model.AppReq) (*model.App, error) {
	db := setupDB()

	var lastInsertID int
	var query string
	var err error
	if app.ServiceId > 0 {
		query = `INSERT INTO apps (name, service_id, allowed_ips, status)
						VALUES ($1, $2, $3, $4) returning id`
		err = db.QueryRow(query, app.Name, app.ServiceId, app.AllowedIPs, STATUS_ACTIVE).Scan(&lastInsertID)
	} else if app.ServiceName != "" {
		query = `INSERT INTO apps (name, service_id, allowed_ips, status)
						VALUES ($1, (SELECT id from services WHERE name=$2 AND domain_id=$3), $4, $5) returning id`
		err = db.QueryRow(query, app.Name, app.ServiceName, domainId, app.AllowedIPs, STATUS_ACTIVE).Scan(&lastInsertID)
	} else {
		err = fmt.Errorf("service must be specified")
	}

	if err != nil {
		return nil, err
	}

	// Select the inserted record and return
	return SelectApp(domainId, "", lastInsertID), err
}

// Select returns the whole database
func SelectApps(domainId int) []*model.App {
	db := setupDB()

	query := `SELECT app.id, app.name, app.service_id, app.allowed_ips, s.name
			FROM apps app, services s LEFT JOIN domains d ON s.domain_id=d.id
			WHERE s.domain_id=$1 AND app.status=$2 AND app.service_id=s.id`
	rows, err := db.Query(query, domainId, STATUS_ACTIVE)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var apps []*model.App

	for {
		app := readAppRow(rows)
		if app == nil {
			break
		}
		apps = append(apps, app)
	}
	return apps
}

// Select the app with either name or id
func SelectApp(domainId int, appName string, appId int) *model.App {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	var app *model.App

	if appId > 0 {
		// Select By App Id
		query = `SELECT app.id, app.name, app.service_id, app.allowed_ips, s.name
				FROM apps app, services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND app.id=$2 AND app.status=$3 AND app.service_id=s.id`
		rows, err = db.Query(query, domainId, appId, STATUS_ACTIVE)
	} else {
		// Select By App Name
		query = `SELECT app.id, app.name, app.service_id, app.allowed_ips, s.name
				FROM apps app, services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND app.name=$2 AND app.status=$3 AND app.service_id=s.id`
		rows, err = db.Query(query, domainId, appName, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("%s: domain=%d [%s %d] %v\n", query, domainId, appName, appId, err)
		return nil
	}
	app = readAppRow(rows)
	if app != nil {
		//log.Printf("Select App %s %d in domain %d - %v\n", appName, appId, domainId, *app)
	} else {
		log.Printf("Select App %s %d in domain %d - NOT FOUND\n", appName, appId, domainId)
	}
	rows.Close()
	return app
}

// Select the app with either service name or service id
func SelectApp2(domainId int, svcName string, svcId int) *model.App {
	db := setupDB()

	var rows *sql.Rows
	var err error
	var query string
	var app *model.App

	if svcId > 0 {
		// Select By App Id
		query = `SELECT app.id, app.name, app.service_id, app.allowed_ips, s.name
				FROM apps app, services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND app.service_id=$2 AND app.status=$3 AND app.service_id=s.id`
		rows, err = db.Query(query, domainId, svcId, STATUS_ACTIVE)
	} else {
		// Select By App Name
		query = `SELECT app.id, app.name, app.service_id, app.allowed_ips, s.name
				FROM apps app, services s LEFT JOIN domains d ON s.domain_id=d.id
				WHERE s.domain_id=$1 AND s.name=$2 AND app.status=$3 AND app.service_id=s.id`
		rows, err = db.Query(query, domainId, svcName, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("%s: domain=%d [%s %d] %v\n", query, domainId, svcName, svcId, err)
		return nil
	}
	app = readAppRow(rows)
	if app != nil {
		//log.Printf("Select App %s %d in domain %d - %v\n", appName, appId, domainId, *app)
	} else {
		log.Printf("Select App %s %d in domain %d - NOT FOUND\n", svcName, svcId, domainId)
	}
	rows.Close()
	return app
}

// Update the the record with the id
func UpdateApp(domainId int, appName string, appId int, app *model.AppReq) *model.App {
	db := setupDB()

	// Compose SQL query
	var params string
	if app.Name != "" {
		params += "name='" + app.Name + "', "
	}
	if app.AllowedIPs != "" {
		params += "allowed_ips='" + app.AllowedIPs + "', "
	}
	if params == "" {
		// Nothing to update
		return SelectApp(domainId, appName, appId)
	}
	params = params[:len(params)-2]

	var err error
	var query string

	// Make sure domain of the service matches domainId specified here
	var rows *sql.Rows
	var did int

	if app.ServiceId > 0 {
		query = "SELECT domain_id FROM services WHERE id=$1"
		rows, _ = db.Query(query, app.ServiceId)
	} else {
		query = "SELECT domain_id FROM services WHERE name=$1"
		rows, _ = db.Query(query, app.ServiceName)
	}
	if !rows.Next() {
		fmt.Printf("No such service with the id or name")
		return nil
	}

	rows.Scan(&did)
	if did != domainId {
		fmt.Printf("The domain of the service doesn't match this domain")
		return nil
	}

	if appId > 0 {
		query = fmt.Sprintf("UPDATE apps SET %s WHERE id=$1 AND status=$2", params)
		_, err = db.Exec(query, appId, STATUS_ACTIVE)
	} else {
		query = fmt.Sprintf("UPDATE apps SET %s WHERE name=$1 AND status=$2", params)
		_, err = db.Exec(query, appName, STATUS_ACTIVE)
	}

	if err != nil {
		fmt.Printf("%s: [%s %d] %v\n", query, appName, appId, err)
		return nil
	}

	// Select the updated record and return
	return SelectApp(domainId, appName, appId)
}

// Delete the record with the id
func DeleteApp(domainId int, appName string, appId int) *model.App {
	db := setupDB()

	deleted_app := SelectApp(domainId, appName, appId)
	if deleted_app == nil {
		// Invalid appName or appId
		fmt.Printf("DeleteApp: [%s %d] domain %d - Invalid App\n", appName, appId, domainId)
		return nil
	}

	var err error
	var query string
	// TODO - Make sure domain of the service matches domainId specified here
	if appId > 0 {
		query = "UPDATE apps SET status=$1 WHERE id=$2 AND status=$3"
		_, err = db.Exec(query, STATUS_DELETED, appId, STATUS_ACTIVE)
	} else {
		query = "UPDATE apps SET status=$1 WHERE name=$2 AND status=$3"
		_, err = db.Exec(query, STATUS_DELETED, appName, STATUS_ACTIVE)
	}
	if err != nil {
		fmt.Printf("DeleteApp: [%s %d] domain %d - %v\n", appName, appId, domainId, err)
		return nil
	}

	return deleted_app
}

func readAppRow(rows *sql.Rows) *model.App {
	var appId int
	var serviceId int
	var name sql.NullString
	var allowedIPs sql.NullString
	var serviceName sql.NullString

	if !rows.Next() {
		return nil
	}

	err := rows.Scan(&appId, &name, &serviceId, &allowedIPs, &serviceName)
	if err != nil {
		fmt.Printf("ReadApp Scan: %v\n", err)
		return nil
	}

	app := model.App{
		ID:          appId,
		Name:        name.String,
		ServiceId:   serviceId,
		ServiceName: serviceName.String,
		AllowedIPs:  allowedIPs.String,
	}
	//fmt.Printf("ReadApp: %s\n", app.Name)
	return &app
}
