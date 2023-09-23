package models

import "encoding/xml"

type CurrencyStorage struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode       int    `xml:"NumCode"`
	CharCode      string `xml:"CharCode"`
	Multiplier    int    `xml:"Nominal"`
	Name          string `xml:"Name"`
	CurrencyValue string `xml:"Value"`
}

type UpdateDatetime struct {
	Id             int    `sql:"id"`
	UpdateDatetime string `sql:"update_datetime"`
}
