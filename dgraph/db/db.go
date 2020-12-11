package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// ErrNotFound indicates requested entity does not exist
var ErrNotFound = fmt.Errorf("no matching record")

type Database struct {
	dg   *dgo.Dgraph
	Conn *grpc.ClientConn
}

func Initialize(addr string) (*Database, error) {
	// establish and check connection
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	d := dgo.NewDgraphClient(
		api.NewDgraphClient(conn),
	)

	log.Print("database connection established")

	db := Database{
		dg:   d,
		Conn: conn,
	}
	err = db.migrateSchema()
	return &db, err
}

func (db Database) migrateSchema() error {
	s := `
	type Item {
		name: string
		description: string
		createdAt: datetime
	}

	name: string .
	description: string .
	createdAt: datetime @index(hour) .
	`
	return db.dg.Alter(context.TODO(), &api.Operation{Schema: s})
}
