package config

type Env struct {
	DB   DBEnv
	Auth Authentication
}

type Authentication struct {
	SigningKey []byte
}
