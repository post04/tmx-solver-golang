package payloadbuilder

import (
	"fmt"

	"github.com/267H/orderedform"
	"github.com/mileusna/useragent"
	"github.com/obfio/tmx-solver-golang/mongo"
)

func browserVer(print *mongo.Print) string {
	agent := useragent.Parse(print.AgentInfo.UserAgent)
	form := orderedform.NewForm(4)
	form.Set("jsou", print.AgentInfo.OSName)
	form.Set("jso", fmt.Sprintf("%s", print.AgentInfo.OSNameAndVersion))
	form.Set("jsbu", print.AgentInfo.BrowserNameAndVersion)
	form.Set("jsb", fmt.Sprintf("%s %v", print.AgentInfo.BrowserNameAndVersion, agent.VersionNo.Major))
	a := "&" + form.URLEncode()
	return a
}
