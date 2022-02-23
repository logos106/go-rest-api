package db

import (
	"fmt"
	"time"
)

func TokenInvalidate(token string) {
	db := setupDB()

	query := `UPDATE sessions SET status=$1, end_time=$2 WHERE session_id=$3 AND status=$4`
	_, err := db.Query(query, STATUS_DELETED, time.Now(), token, STATUS_ACTIVE)
	if err != nil {
		fmt.Printf("TokenInvalidate: %s %v\n", token, err)
	} else {
		fmt.Printf("TokenInvalidate: %s SUCCESS\n", token)
	}
}
