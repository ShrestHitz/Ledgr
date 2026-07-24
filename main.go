package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/shresth/ledgr/db/sqlc"
	"github.com/shresth/ledgr/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connPool)
	_ = store

	log.Println("Ledgr server starting on", config.ServerAddress)

	// TODO: wire up Gin router and start HTTP server
}
