package db

import (
	"database/sql"
	"fmt"

	m "github.com/saroopmathur/rest-api/models"
)

func GetUserPolicy(domainId int, userName string, userId int) (*m.Policy, error) {
	db := setupDB()

	var err error
	policy := &m.Policy{}
	policy.ServiceNodes = make(map[string]*m.ServiceNode)

	// Get userId from userName (Usually if a username is specified by Admin)
	u := SelectUser(domainId, userName, userId)
	if u == nil {
		// Bad username (since userId is 0)
		err = fmt.Errorf("username '%s' invalid", userName)
		return nil, err
	}
	userId = u.ID

	// Get user specific policies
	query := `SELECT services.name, services.wg_key, services.virtual_ip, services.public_ip, services.local_ip,
			apps.name, apps.apps
				FROM services, apps, user_access_control ua
				WHERE ua.user_id=$1
					AND ua.app_id=apps.id
					AND apps.service_id=services.id`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	readPolicyRows(rows, policy)
	rows.Close()

	for _, s := range policy.ServiceNodes {
		for _, app := range s.Apps {
			app.IsUserPolicy = true
		}
	}

	// Get group policies
	query = `SELECT services.name, services.wg_key, services.virtual_ip, services.public_ip, services.local_ip,
			apps.name, apps.apps
				FROM services, apps, group_access_control ga
				WHERE ga.group_id IN (SELECT DISTINCT members.group_id FROM group_members members
							WHERE members.user_id=$1)
					AND ga.app_id=apps.id
					AND apps.service_id=services.id`

	rows, err = db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	readPolicyRows(rows, policy)
	rows.Close()
	return policy, nil
}

func GetAllPolicies(domainId int) (*m.Policy, error) {
	db := setupDB()

	var err error
	policy := &m.Policy{}
	policy.ServiceNodes = make(map[string]*m.ServiceNode)

	// Get user specific policies
	query := `SELECT services.name, services.wg_key, services.virtual_ip, services.public_ip, services.local_ip,
			apps.name, apps.apps
				FROM services, apps, user_access_control ua
				WHERE ua.app_id=apps.id
					AND apps.service_id=services.id
					AND apps.status=$1
					AND services.status=$1
					AND services.domain_id=$2`

	rows, err := db.Query(query, STATUS_ACTIVE, domainId)
	if err != nil {
		return nil, err
	}
	readPolicyRows(rows, policy)
	rows.Close()

	for _, s := range policy.ServiceNodes {
		for _, app := range s.Apps {
			app.IsUserPolicy = true
		}
	}

	// Get group policies
	query = `SELECT services.name, services.wg_key, services.virtual_ip, services.public_ip, services.local_ip,
			apps.name, apps.apps
				FROM services, apps, group_access_control ga
				WHERE ga.app_id=apps.id
					AND apps.service_id=services.id
					AND apps.status=$1
					AND services.status=$1
					AND services.domain_id=$2`

	rows, err = db.Query(query, STATUS_ACTIVE, domainId)
	if err != nil {
		return nil, err
	}
	readPolicyRows(rows, policy)
	rows.Close()
	return policy, nil
}

func readPolicyRows(rows *sql.Rows, policy *m.Policy) {
	var serviceName sql.NullString
	var wgKey sql.NullString
	var vip sql.NullString
	var public_ip sql.NullString
	var local_ip sql.NullString
	var appName sql.NullString
	var allowedIPs sql.NullString

	if rows.Next() {
		err := rows.Scan(&serviceName, &wgKey, &vip, &public_ip, &local_ip, &appName, &allowedIPs)
		checkErr(err)

		app := &m.PolicyApp{}
		app.Name = appName.String
		app.AllowedIPs = allowedIPs.String

		service := policy.ServiceNodes[serviceName.String]
		if service == nil {
			// New service
			service = &m.ServiceNode{}
			service.Name = serviceName.String
			service.WGKey = wgKey.String
			service.VirtualIP = vip.String
			service.PublicIP = public_ip.String
			service.LocalIP = local_ip.String
			policy.ServiceNodes[service.Name] = service
		}
		service.Apps = append(service.Apps, app)
	}
}
