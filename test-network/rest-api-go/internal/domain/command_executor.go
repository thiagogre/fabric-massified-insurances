package domain

type CommandExecutorInterface interface {
	ExecuteCommand(name string, args ...string) ([]byte, error)
}
