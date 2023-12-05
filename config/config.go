package config

type AuthenticationConfig struct {
	HTTPAddress struct {
		Port string `yaml:"port" envconfig:"HTTP_PORT"`
		Host string `yaml:"host" envconfig:"HTTP_HOST"`
	} `yaml:"http"`
	GRPCAddress struct {
		Port string `yaml:"port" envconfig:"GRPC_PORT"`
		Host string `yaml:"host" envconfig:"GRPC_HOST"`
	} `yaml:"grpc"`
	DbConnection struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Database string `yaml:"database" envconfig:"DB_DATABASE"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	} `yaml:"database"`
	SecretKey string `yaml:"secret" envconfig:"SECRET_KEY"`
}

type WatermarkConfig struct {
	HTTPAddress struct {
		Port string `yaml:"port" envconfig:"HTTP_PORT"`
		Host string `yaml:"host" envconfig:"HTTP_HOST"`
	} `yaml:"http"`
	GRPCAddress struct {
		Port string `yaml:"port" envconfig:"GRPC_PORT"`
		Host string `yaml:"host" envconfig:"GRPC_HOST"`
	} `yaml:"grpc"`
	DbConnection struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Database string `yaml:"database" envconfig:"DB_DATABASE"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	} `yaml:"database"`
	Cloudinary struct {
		Cloud  string `yaml:"cloud" envconfig:"CLOUDINARY_CLOUD"`
		Api    string `yaml:"api" envconfig:"CLOUDINARY_API"`
		Secret string `yaml:"secret" envconfig:"CLOUDINARY_SECRET"`
	} `yaml:"cloudinary"`
	Services struct {
		Auth struct {
			Port string `yaml:"port" envconfig:"AUTH_PORT"`
			Host string `yaml:"host" envconfig:"AUTH_HOST"`
		} `yaml:"auth"`
		Picture struct {
			Port string `yaml:"port" envconfig:"PICTURE_PORT"`
			Host string `yaml:"host" envconfig:"PICTURE_HOST"`
		} `yaml:"picture"`
	} `yaml:"services"`
}

type PictureConfig struct {
	HTTPAddress struct {
		Port string `yaml:"port" envconfig:"HTTP_PORT"`
		Host string `yaml:"host" envconfig:"HTTP_HOST"`
	} `yaml:"http"`
	GRPCAddress struct {
		Port string `yaml:"port" envconfig:"GRPC_PORT"`
		Host string `yaml:"host" envconfig:"GRPC_HOST"`
	} `yaml:"grpc"`
}
