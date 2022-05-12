package main

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var gml string = `
	<fme:GML gml:id="id5255fa48-42b3-43d1-9e0d-b3ba8b57a936">
		<fme:OBJECTID>1</fme:OBJECTID>
		<fme:RECORD_ID>TEST48234</fme:RECORD_ID>
		<fme:NAME>HATCHPIN POINT</fme:NAME>
		<fme:FEAT_CODE>PT</fme:FEAT_CODE>
		<fme:CGDN>N</fme:CGDN>
		<fme:AUTHORITY_ID>OLD</fme:AUTHORITY_ID>
		<fme:CONCISE_GAZ>N</fme:CONCISE_GAZ>
		<fme:LATITUDE>-12.58361</fme:LATITUDE>
		<fme:lat_degrees>-12</fme:lat_degrees>
		<fme:lat_minutes>35</fme:lat_minutes>
		<fme:lat_seconds>0</fme:lat_seconds>
		<fme:LONGITUDE>141.62583</fme:LONGITUDE>
		<fme:long_degrees>141</fme:long_degrees>
		<fme:long_minutes>37</fme:long_minutes>
		<fme:long_seconds>32</fme:long_seconds>
		<fme:STATE_ID>OLD</fme:STATE_ID>
		<fme:STATUS>U</fme:STATUS>
		<fme:VARIANT_NAME/>
		<fme:MAP_100K>7272</fme:MAP_100K>
		<fme:Place_ID>45880</fme:Place_ID>
		<gml:pointProperty>
			<gml:Point srsName="EPSG:3112" srsDimension="2">
				<gml:pos>141.625915527344 -12.5836181640625</gml:pos>
			</gml:Point>
		</gml:pointProperty>
	</fme:GML>
`

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestMain(t *testing.T) {
	b := []byte(gml)
	var FeatureMember FeatureMember
	err := xml.Unmarshal(b, &FeatureMember)
	if err != nil {
		panic(err)
	}
	fmt.Println(FeatureMember)

	equals(t, "HATCHPIN POINT", FeatureMember.Name)
	equals(t, "PT", FeatureMember.FCode)
	equals(t, "-12.58361", FeatureMember.DDLat)
	equals(t, "141.62583", FeatureMember.DDLon)
	equals(t, "", FeatureMember.Variants)
	equals(t, "OLD", FeatureMember.State)
	equals(t, "EPSG:3112", FeatureMember.PP.Pt.EPSG)
	equals(t, "U", FeatureMember.Status)
	equals(t, "-12", FeatureMember.LatDeg)
	equals(t, "35", FeatureMember.LatMin)
	equals(t, "0", FeatureMember.LatSec)
	equals(t, "141", FeatureMember.LonDeg)
	equals(t, "37", FeatureMember.LonMin)
	equals(t, "32", FeatureMember.LonSec)

}
