package settings

//Window type
type Window struct {
	Width  float64
	Height float64
}

//Window size
func GetWindowDefault() Window {
	//return Window{1280, 720} //HD resolution
	return Window{1600, 900} //HD+ resolution
	//return Window{1920, 1080}//Full HD resolution
}
