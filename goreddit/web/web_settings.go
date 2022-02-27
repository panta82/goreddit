package web

type Settings struct {
	Port int
}

func NewSettings() *Settings {
	return &Settings{
		Port: 3000,
	}
}
