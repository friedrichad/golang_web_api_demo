package mail

import (
	"github.com/spf13/viper"
)
type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}
var Mail MailConfig

type Config struct {
	Mail MailConfig `mapstructure:"mail"`
}


func InitMailConfig() {
	Mail = MailConfig{
		Host: viper.GetString("mail.host"),
		Port: viper.GetInt("mail.port"),
		Username: viper.GetString("mail.username"),
		Password: viper.GetString("mail.password"),
		From: viper.GetString("mail.from"),
	}
}