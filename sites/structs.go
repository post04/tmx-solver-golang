package sites

type Profile struct {
	TimezoneOffset          int
	SecondaryTimezoneOffset int
	ScreenWidth             int
	ScreenHeight            int
	ScreenAnvilWidth        int
	ScreenAnvilHeight       int
	InnerWidth              int
	InnerHeight             int
	OuterWidth              int
	OuterHeight             int
	MimeTypes               string
	Plugins                 string
	Math                    string
	PluginSupport           string

	MimeTypesLength int
	PluginsLength   int

	ColorDepth int

	HardwareConcurrency int
	DeviceMemory        int
	MaxTouchPoints      int

	Timezone string
}

type UserAgent struct {
	UA             string
	BrowserName    string
	BrowserVersion string
	OSName         string
	OSVersion      string
}
