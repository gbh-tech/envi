package providers

type EnvProvider interface {
	GenerateEnvFile() error
}
