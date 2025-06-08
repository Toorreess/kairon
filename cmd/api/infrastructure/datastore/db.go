package datastore

import (
	"context"
	db "kairon/adapters/database"
	"kairon/adapters/database/clients/firestore"
	"kairon/config"
	"fmt"
)

func NewDBConnection() (*db.Connection, error) {
	var conn db.Connection
	conn.Ctx = context.Background()

	switch config.C.Database.DBType {
	case "firestore":
		client, err := firestore.NewFirestoreClient(config.C.ProjectID)
		if err != nil {
			return nil, err
		}
		conn.Client = client
	default:
		return &conn, fmt.Errorf("Invalid DB type")
	}

	return &conn, nil
}
