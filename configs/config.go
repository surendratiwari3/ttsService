package configs

import (
	"github.com/spf13/viper"
	"strings"
)

type awsConfig struct {
	AccessKeyID    string
	SecretAcessKey string
	Region         string
}

type echoConfig struct {
	HostPort string
}

type minioConfig struct {
	Host		string
	AccessKeyID    string
	SecretAcessKey string
	Secure 	bool
}

// Config - configuration object
type Config struct {
	AwsConfig   awsConfig
	EchoConfig  echoConfig
	MinioConfig minioConfig
}

var conf *Config

// GetConfig - Function to get Config
func GetConfig() *Config {
	if conf != nil {
		return conf
	}
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	awsConf := awsConfig{
		AccessKeyID:    v.GetString("aws.id"),
		SecretAcessKey: v.GetString("aws.key"),
		Region:         v.GetString("aws.region"),
	}

	minioConf := minioConfig{
		Host:v.GetString("minio.host"),
		AccessKeyID:    v.GetString("minio.id"),
		SecretAcessKey: v.GetString("minio.key"),
		Secure:v.GetBool("minio.secure"),
	}

	echoConf := echoConfig{
		HostPort: v.GetString("echo.hostport"),
	}
	conf = &Config{
		AwsConfig:   awsConf,
		EchoConfig:  echoConf,
		MinioConfig: minioConf,
	}
	return conf
}
