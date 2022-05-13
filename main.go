package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func main() {
	importFileName := "example.gml"
	importFile, iferr := os.Open(importFileName)
	checkError(iferr, "iferror")
	fmt.Println("Gazetteer file open")

	defer importFile.Close()

	exportFileName := "gazetteer.csv"
	outfile, oferr := os.Create(exportFileName)
	checkError(oferr, "oferror")

	defer outfile.Close()

	csvHeaders := []string{"Name", "Variants", "State", "Decimal latitude",
		"Decimal longitude", "Latitude degrees", "Latitude minutes", "Latitude seconds",
		"Longitude degrees", "Longitude minutes", "Longitude seconds", "Feature Code", "Status",
		"Datum", "EPSG"}
	export := csv.NewWriter(outfile)
	export.Write(csvHeaders)

	dec := xml.NewDecoder(importFile)

	for {
		t, tokenErr := dec.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			checkError(tokenErr, "tokenerror")
		}
		switch ty := t.(type) {
		case xml.StartElement:
			if ty.Name.Local == "GML" {
				var FeatureMember FeatureMember
				dec.DecodeElement(&FeatureMember, &ty)
				FeatureMember.EPSG = FeatureMember.PP.Pt.EPSG
				FeatureMember.writeToCSV(*export)
			}
		}
	}
}

type Point struct {
	EPSG string `xml:"srsName,attr"`
}

type PointProperty struct {
	Pt Point `xml:"Point"`
}

type FeatureMember struct {
	Name     string        `xml:"NAME"`
	Variants string        `xml:"VARIANT_NAME"`
	State    string        `xml:"STATE_ID"`
	DDLat    string        `xml:"LATITUDE"`
	DDLon    string        `xml:"LONGITUDE"`
	LatDeg   string        `xml:"lat_degrees"`
	LatMin   string        `xml:"lat_minutes"`
	LatSec   string        `xml:"lat_seconds"`
	LonDeg   string        `xml:"long_degrees"`
	LonMin   string        `xml:"long_minutes"`
	LonSec   string        `xml:"long_seconds"`
	FCode    string        `xml:"FEAT_CODE"`
	Status   string        `xml:"STATUS"`
	PP       PointProperty `xml:"pointProperty"`
	EPSG     string
	Datum    string
}

func (fm *FeatureMember) writeToCSV(exportFile csv.Writer) {
	var locData [15]string
	locData[0] = fm.Name
	locData[1] = fm.Variants
	locData[2] = fm.State
	locData[3] = fm.DDLat
	locData[4] = fm.DDLon
	locData[5] = fm.LatDeg
	locData[6] = fm.LatMin
	locData[7] = fm.LatSec
	locData[8] = fm.LonDeg
	locData[9] = fm.LonMin
	locData[10] = fm.LonSec
	locData[11] = fm.FCode
	locData[12] = fm.Status
	locData[13] = "GDA94"
	locData[14] = fm.EPSG

	exportFile.Write(locData[:])
	exportFile.Flush()
}

func checkError(err error, code string) {
	if err != nil {
		switch code {
		case "iferror":
			fmt.Println("ERROR: Could not open gazetteer import file (Expecting Gazetteer2012GML.gml)")
		case "oferror":
			fmt.Println("ERROR: Could not create gazetteer export file (Attempting to create gazetteer.csv)")
		case "tokenerror":
			fmt.Println("ERROR: Did not obtain expected XML token")
			panic(err)
		}
		os.Exit(1)
	}
}
