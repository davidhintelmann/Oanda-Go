package restful

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/sharedmemory"
)

func ConnectMSSQL(
	ctx context.Context,
	conn *sql.DB,
	driver string,
	server string,
	database string,
	trusted_connection bool,
	encrypt bool) (*sql.DB, error) {
	var err error

	connString := fmt.Sprintf("server=%s;database=%s;TrustServerCertificate=%v;encrypt=%v", server, database, trusted_connection, encrypt)
	conn, err = sql.Open(driver, connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
		return nil, err
	}

	// log.Printf("Connected!\n")

	// Ping database to see if it's still alive.
	// Important for handling network issues and long queries.
	// err = db.PingContext(ctx)
	// if err != nil {
	// 	log.Fatal("Error pinging database: " + err.Error())
	// }
	return conn, nil
}
