package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	server "test-backend-1-curboturbo/internal/server"
)


func main(){
	errorChan := make(chan os.Signal,1)
	signal.Notify(errorChan, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
	storagePath := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
    	os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("POSTGRES_DB"),
    )
	db, _ := gorm.Open(postgres.Open(storagePath), &gorm.Config{})
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
        panic("failed to enable pgcrypto:")
    }
	
	server := server.New(db)
	go func(){
		if err :=server.Run();err!=nil{
			errorChan<-os.Interrupt
		}
	}()
	<-errorChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Stop(ctx)

}