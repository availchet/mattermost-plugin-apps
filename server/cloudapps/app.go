package cloudapps

type AppID string

type App struct {
	AppID
	DisplayName string
	Description string
	RootURL     string
}
