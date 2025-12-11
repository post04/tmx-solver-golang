package payloadbuilder

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/267H/orderedform"
	"github.com/mileusna/useragent"
)

func browserGeneral(r *Request) string {
	timezoneOffset, _, DSTShift, err := getDSTShift(r.Print.TimeZone.SecondaryTimeZoneOffset)
	if err != nil {
		timezoneOffset = -600
		DSTShift = 0
	}
	agent := useragent.Parse(r.Print.AgentInfo.UserAgent)
	form := orderedform.NewForm(25)
	form.Set("c", fmt.Sprint(timezoneOffset))
	form.Set("z", fmt.Sprint(DSTShift))
	form.Set("f", fmt.Sprintf(`%vx%v`, r.Print.Screen.Width, r.Print.Screen.Height))
	form.Set("af", fmt.Sprintf(`%vx%v`, r.Print.Screen.OuterWidth, r.Print.Screen.OuterHeight))
	form.Set("sxy", fmt.Sprintf("%vx%v", r.Print.Screen.ScreenX, r.Print.Screen.ScreenY))
	form.Set("dpr", fmt.Sprintf("1,%d,%d,%d,%d,%d,%d,%d,%d,0,0", r.Print.Screen.Width, r.Print.Screen.Height, r.Print.Screen.OuterWidth, r.Print.Screen.OuterHeight, r.Print.Screen.InnerWidth, r.Print.Screen.InnerHeight, r.Print.Screen.OuterWidth, r.Print.Screen.OuterHeight))
	form.Set("mt", md5Hex(r.Print.MimeTypes.MimeTypesStr))
	form.Set("mn", fmt.Sprintf("%d", r.Print.MimeTypes.Length))
	form.Set("scd", fmt.Sprintf("%d", r.Print.Screen.ColorDepth))
	// LH should be no more than 255 characters
	lh := r.URL
	if len(lh) > 255 {
		lh = lh[:255]
	}
	form.Set("lh", lh)
	// lh -> dr -> p
	// seems everything else is removed right now?
	form.Set("pl", fmt.Sprintf("%d", r.Print.Plugins.Length))
	form.Set("ph", md5Hex(r.Print.Plugins.PluginsStr))
	form.Set("hh", md5Hex(r.OrgID+strings.ToLower(r.SessionID)))
	form.Set("jso", r.Print.AgentInfo.OSNameAndVersion)
	form.Set("jsb", r.Print.AgentInfo.BrowserNameAndVersion+" "+fmt.Sprint(agent.VersionNo.Major))
	form.Set("jsou", r.Print.AgentInfo.OSName)
	form.Set("jsbu", r.Print.AgentInfo.BrowserNameAndVersion)
	form.Set("nhc", fmt.Sprintf("%d", r.Print.General.HardwareConcurrency))
	// TODO: get device memory
	form.Set("ndm", fmt.Sprintf("%d", 8))
	form.Set("nmtp", fmt.Sprintf("%d", r.Print.Screen.TouchPoints))
	form.Set("tzd", r.Print.TimeZone.SecondaryTimeZoneOffset)
	form.Set("mathr", sha256Hex(r.Print.General.MathStuff))
	// urlParts := strings.Split(r.URL, "/")
	// if len(urlParts) < 3 {
	// 	form.Set("dr", "https://www.canadiantire.ca/")
	// } else {
	// 	form.Set("dr", "https://"+urlParts[2])
	// }
	parts := strings.Split(r.URL, "?")
	if len(parts) > 1 {
		form.Set("dr", parts[0])
	} else {
		form.Set("dr", r.URL)
	}
	form.Set("p", r.Print.General.PluginSupport)
	form.Set("ccd", fmt.Sprint(3+rand.Intn(10)))
	a := "&" + form.URLEncode()
	return a
}

func getOffsetMinutes(t time.Time) int {
	_, offset := t.Zone()
	return offset / 60 // convert seconds to minutes
}

func getDSTShift(tzName string) (minOffset, maxOffset, dstShift int, err error) {
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return 0, 0, 0, err
	}

	// Use the same year to compare summer vs winter offset
	year := 2025

	// June 1st (DST in many regions)
	june := time.Date(year, time.June, 1, 0, 0, 0, 0, loc)
	juneOffset := getOffsetMinutes(june)

	// December 1st (Standard Time)
	dec := time.Date(year, time.December, 1, 0, 0, 0, 0, loc)
	decOffset := getOffsetMinutes(dec)

	if juneOffset < decOffset {
		minOffset, maxOffset = juneOffset, decOffset
	} else {
		minOffset, maxOffset = decOffset, juneOffset
	}
	dstShift = maxOffset - minOffset

	return minOffset, maxOffset, dstShift, nil
}
