package plugin

type Plugin interface {
	Init() error
	Start() error
	Stop() error
}
