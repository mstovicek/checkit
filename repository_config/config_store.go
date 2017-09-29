package repository_config

type ConfigStore interface {
	HasConfig(server string, repositoryName string) bool
	LoadConfig(server string, repositoryName string) (Config, error)
	SaveConfig(config Config) error
}
