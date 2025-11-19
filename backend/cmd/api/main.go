package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/server"
)


func main() {
    cfg := config.Load()

    // Подключение к базе данных
    db, err := database.NewDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

	srv := server.NewServer(db)
    
    log.Fatal(srv.Start("8080"))

}
