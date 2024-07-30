package resource

//go:generate ./fyne bundle -package resource -o logo.go assets

var (
	AppLogo = resourceLogoPng
)
