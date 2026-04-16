package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"test-backend-1-curboturbo/internal/service"
	"time"
	"gorm.io/gorm"
	router "test-backend-1-curboturbo/internal/adapters/inbound/gin"
	conferenceAPIAdapter "test-backend-1-curboturbo/internal/adapters/outbound/API/conference"
	adapterStorage "test-backend-1-curboturbo/internal/adapters/outbound/storage"
	adapterJWT "test-backend-1-curboturbo/internal/adapters/outbound/tokenizer"
)


type Server struct{
	http *http.Server
}

func New(db *gorm.DB) *Server{
	token_adapter := adapterJWT.NewTokenGenerator()
	storage_adapter := adapterStorage.NewStorage(db)
	conference_adapter := conferenceAPIAdapter.NewLinkConferenceService()

	auth_service := service.NewAuthService(storage_adapter,token_adapter)
	booking_service := service.NewRoomService(storage_adapter,conference_adapter)
	
	router := router.RouterInit(auth_service,booking_service, token_adapter)

	port := os.Getenv("HTTP_PORT")
	timeoutStr := os.Getenv("TIMEOUT")
	timeoutVal, _ := strconv.Atoi(timeoutStr)
	serv := &Server{
		http: &http.Server{
			Addr: ":"+port,
			ReadTimeout: time.Duration(timeoutVal) * time.Second,
			WriteTimeout: time.Duration(timeoutVal) * time.Second,
			Handler: router,
		},
	}
	return serv
}

func (s *Server) Run() (error){
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) (error){
	return s.http.Shutdown(ctx)
}