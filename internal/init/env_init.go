package init
import (
	configs "test-backend-1-curboturbo/config"
	"github.com/ilyakaznacheev/cleanenv"
	"fmt"
)


func LoadNewConfig(path string) (*configs.AppConfig, error, string){
	var cfg configs.AppConfig
	if err := cleanenv.ReadConfig(path, &cfg);err!=nil{
		return nil, err, ""
	}
	connectionPath := CreateStorage(cfg.DataBase)
	return &cfg, nil, connectionPath
}

func CreateStorage(cfg configs.PostgreConfig) string{
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
    )
}