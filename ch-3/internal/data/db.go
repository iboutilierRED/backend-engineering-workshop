package data

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

func InitDbConnection() (gocqlx.Session, error) {
	cluster := gocql.NewCluster("localhost:9042")

	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Canhassi", Password: "password123"}
	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("AWS_US_EAST_1")

	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		// cannot return nil for type gocqlx.Session, presence of error will suffice
		return gocqlx.Session{}, err
	}
	log.Println("Database connected")
	createKeyspace := `CREATE KEYSPACE IF NOT EXISTS wikistats 
					   WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': '3'}  
					   AND durable_writes = true;`

	if err = session.Query(createKeyspace, nil).ExecRelease(); err != nil {
		return gocqlx.Session{}, err
	}

	// Now that keyspace has been created, we need to set it and re-connect
	cluster.Keyspace = "wikistats"
	session, err = gocqlx.WrapSession(cluster.CreateSession())

	return session, nil
}
