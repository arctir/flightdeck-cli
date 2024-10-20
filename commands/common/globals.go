package common

type Globals struct {
	APIEndpoint  string `name:"endpoint" env:"ARCTIR_API_ENDPOINT" default:"${apiEndpoint}"`
	ConfigPath   string `name:"config" env:"ARCTIR_CONFIG" default:"${configPath}"`
	OutputFormat string `name:"format" enum:"json,table" default:"table"`

	Debug bool
}
