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

type WatermarkConfig struct {
	HTTPAddress  address            `yaml:"http"`
	GRPCAddress  address            `yaml:"grpc"`
	DbConnection DatabaseConnection `yaml:"database"`
	Cloudinary   struct {
		Cloud  string `yaml:"cloud"`
		Api    string `yaml:"api"`
		Secret string `yaml:"secret"`
	} `yaml:"cloudinary"`
	Services struct {
		Auth    address `yaml:"auth"`
		Picture address `yaml:"picture"`
	} `yaml:"services"`
}

type PictureConfig struct {
	HTTPAddress  address            `yaml:"http"`
	GRPCAddress  address            `yaml:"grpc"`
	DbConnection DatabaseConnection `yaml:"database"`
}
