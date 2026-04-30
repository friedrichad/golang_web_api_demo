package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"log"
)

func InitMinioClient() *minio.Client {
	client, err := minio.New(viper.GetString("minio.url"), &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("minio.access-key"), viper.GetString("minio.secret-key"), ""),
		Secure: viper.GetBool("minio.secure"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
