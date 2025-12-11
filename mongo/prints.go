package mongo

import (
	"math/rand"

	"github.com/obfio/tmx-solver-golang/config"
)

type Print struct {
	AgentInfo struct {
		OSName                string `json:"OSName"`
		BrowserName           string `json:"BrowserName"`
		BrowserNameAndVersion string `json:"BrowserNameAndVersion"`
		OSNameAndVersion      string `json:"OSNameAndVersion"`
		UserAgent             string `json:"UserAgent"`
		Vendor                string `json:"Vendor"`
		Platform              string `json:"Platform"`
		AppVersion            string `json:"AppVersion"`
	} `json:"agentInfo"`
	TimeZone struct {
		TimeZoneOffset          int    `json:"TimeZoneOffset"`
		SecondaryTimeZoneOffset string `json:"SecondaryTimeZoneOffset"`
	} `json:"TimeZone"`
	Screen struct {
		Width           int  `json:"Width"`
		Height          int  `json:"Height"`
		AvailableWidth  int  `json:"AvailableWidth"`
		AvailableHeight int  `json:"AvailableHeight"`
		ScreenX         int  `json:"ScreenX"`
		ScreenY         int  `json:"ScreenY"`
		InnerWidth      int  `json:"InnerWidth"`
		InnerHeight     int  `json:"InnerHeight"`
		OuterWidth      int  `json:"OuterWidth"`
		OuterHeight     int  `json:"OuterHeight"`
		ColorDepth      int  `json:"ColorDepth"`
		IsTouchCapable  bool `json:"IsTouchCapable"`
		TouchPoints     int  `json:"TouchPoints"`
	} `json:"Screen"`
	MimeTypes struct {
		Length int `json:"Length"`
		// MimeTypes []struct {
		// 	Type        string `json:"Type"`
		// 	Description string `json:"Description"`
		// 	Suffixs     string `json:"Suffixs"`
		// } `json:"MimeTypes"`
		MimeTypesStr string `json:"MimeTypesStr"`
	} `json:"MimeTypes"`
	Plugins struct {
		Length int `json:"Length"`
		// Plugins []struct {
		// 	Name        string `json:"Name"`
		// 	Description string `json:"Description"`
		// 	Filename    string `json:"Filename"`
		// 	Length      int    `json:"Length"`
		// } `json:"Plugins"`
		PluginsStr string `json:"PluginsStr"`
	} `json:"Plugins"`
	General struct {
		Language            string `json:"Language"`
		HardwareConcurrency int    `json:"HardwareConcurrency"`
		Timezone            string `json:"Timezone"`
		MathStuff           string `json:"MathStuff"`
		PluginSupport       string `json:"PluginSupport"`
		// Battery             any    `json:"Battery"`
		BatteryStr          string `json:"BatteryStr"`
		UserAgentDump       string `json:"UserAgentDump"`
		UserAgentBrandsDump string `json:"UserAgentBrandsDump"`
		MediaDevicesStr     string `json:"MediaDevicesStr"`
	} `json:"General"`
	Wgl struct {
		// WGLBlob1      string `json:"WGLBlob1"`
		WGLHash1 string `json:"WGLHash1"`
		// WGLBlob2      string `json:"WGLBlob2"`
		WGLHash2      string `json:"WGLHash2"`
		WglC          string `json:"WGL_C"`
		WglAgentHash  string `json:"WglAgentHash"`
		WGLVendor     string `json:"WGLVendor"`
		WGLRenderer   string `json:"WGLRenderer"`
		WGLAgents     string `json:"WGLAgents"`
		WglAgentHash2 string `json:"WglAgentHash2"`
	} `json:"WGL"`
	Canvas struct {
		// CanvasBlob string `json:"CanvasBlob"`
		CanvasHash string `json:"CanvasHash"`
	} `json:"Canvas"`
	Audio struct {
		// AudioBlob string `json:"AudioBlob"`
		AudioHash string `json:"AudioHash"`
	} `json:"Audio"`
	Fonts struct {
		Count int    `json:"count"`
		Str   string `json:"str"`
	} `json:"Fonts"`
}

var printsColl = config.DbCnx.Collection("prints")

var (
	prints []*Print
)

func GetRandomPrint() *Print {
	return prints[rand.Intn(len(prints))]
}
