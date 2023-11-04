package config

type DatabaseConnection struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
}
type address struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type AuthenticationConfig struct {
	HTTPAddress  address            `yaml:"http"`
	GRPCAddress  address            `yaml:"grpc"`
	DbConnection DatabaseConnection `yaml:"database"`
	SecretKey    string             `yaml:"secret"`
}

type DatabaseConfig struct {
	HTTPAddress  address            `yaml:"http"`
	GRPCAddress  address            `yaml:"grpc"`
	DbConnection DatabaseConnection `yaml:"database"`
	Services     struct {
		Auth address `yaml:"auth"`
	} `yaml:"services"`
}

type WatermarkConfig struct {
	HTTPAddress  address            `yaml:"http"`
	GRPCAddress  address            `yaml:"grpc"`
	DbConnection DatabaseConnection `yaml:"database"`
}
