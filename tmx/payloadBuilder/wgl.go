package payloadbuilder

import (
	"github.com/267H/orderedform"
	"github.com/obfio/tmx-solver-golang/mongo"
)

func wgl(p *mongo.Print) string {
	form := orderedform.NewForm(13)
	// device battery
	form.Set("batst", `{"level":1.00,"status":"charging"}`)

	// audio hash, unsure how to emulate this yet
	// TODO: emulate this
	form.Set("audh", p.Audio.AudioHash)

	// canvas/webgl hashes
	form.Set("ex3", p.Canvas.CanvasHash)
	form.Set("ex4", p.Wgl.WGLHash1)
	form.Set("ex5", p.Wgl.WGLHash2)

	// webgl agent
	form.Set("gl_c", p.Wgl.WglC)
	form.Set("gl_h", p.Wgl.WglAgentHash)
	form.Set("wglv", p.Wgl.WGLVendor)
	form.Set("wglr", p.Wgl.WGLRenderer)
	form.Set("glh_h", p.Wgl.WglAgentHash2)

	// operating system + version
	form.Set("jso", p.AgentInfo.OSNameAndVersion)

	// user agent print
	form.Set("uah", p.General.UserAgentDump)

	// user agent brands print
	form.Set("ual", p.General.UserAgentBrandsDump)
	a := "&" + form.URLEncode()
	return a
}
