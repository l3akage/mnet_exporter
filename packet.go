package main

import "encoding/xml"

type Packet struct {
	XMLName         xml.Name        `xml:"Packet"`
	Command         string          `xml:"Command"`
	DatabaseManager DatabaseManager `xml:"DatabaseManager"`
}

type DatabaseManager struct {
	Mnet []Mnet `xml:"Mnet"`
}

type Mnet struct {
	Group          string `xml:"Group,attr"`
	Bulk           string `xml:"Bulk,attr"`
	EnergyControl  string `xml:"EnergyControl,attr"`
	SetbackControl string `xml:"SetbackControl,attr"`
	ScheduleAvail  string `xml:"ScheduleAvail,attr"`
}
