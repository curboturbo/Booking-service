package server

import (
	"context"
	"net/http"
	"strconv"
	config "test-backend-1-curboturbo/config"
	router "test-backend-1-curboturbo/internal/adapters/inbound/gin"
	adapterStorage "test-backend-1-curboturbo/internal/adapters/outbound/storage"
	adapterJWT "test-backend-1-curboturbo/internal/adapters/outbound/tokenizer"
	conferenceAPIAdapter "test-backend-1-curboturbo/internal/adapters/outbound/API/conference"
	"test-backend-1-curboturbo/internal/service"

	//"github.com/vertica/vertica-sql-go/logger"
	"gorm.io/gorm"
)



type Server struct{
	cfg *config.AppConfig
	http *http.Server
	//logger *log.Logger
}

func New(cfg *config.AppConfig, db *gorm.DB) *Server{
	token_adapter := adapterJWT.NewTokenGenerator()
	storage_adapter := adapterStorage.NewStorage(db)
	conference_adapter := conferenceAPIAdapter.NewLinkConferenceService()

	auth_service := service.NewAuthService(storage_adapter,token_adapter)
	booking_service := service.NewRoomService(storage_adapter,conference_adapter)
	
	router := router.RouterInit(auth_service,booking_service, token_adapter)

	serv := &Server{
		cfg: cfg,
		http: &http.Server{
			Addr: ":"+strconv.Itoa(cfg.Server.Port),
			ReadTimeout: cfg.Server.TimeOut,
			WriteTimeout: cfg.Server.TimeOut,
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