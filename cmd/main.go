package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	loader "test-backend-1-curboturbo/internal/init"
	server "test-backend-1-curboturbo/internal/server"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func main(){
	errorChan := make(chan os.Signal,1)
	signal.Notify(errorChan, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
	cfg, err, storagePath := loader.LoadNewConfig("../config/config.yaml")
	if err != nil{
		panic("cannot load configs")
	}
	db, err := gorm.Open(postgres.Open(storagePath), &gorm.Config{})
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
        panic("failed to enable pgcrypto:")
    }
	server := server.New(cfg, db)
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