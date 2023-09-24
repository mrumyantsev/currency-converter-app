package models

import "encoding/xml"

type CurrencyStorage struct {
	XMLName    xml.Name   `xml:"ValCurs" json:"currencyStorage"`
	Currencies []Currency `xml:"Valute" json:"currencies"`
}

type Currency struct {
	NumCode       int    `xml:"NumCode" json:"numCode"`
	CharCode      string `xml:"CharCode" json:"charCode"`
	Multiplier    int    `xml:"Nominal" json:"multiplier"`
	Name          string `xml:"Name" json:"name"`
	CurrencyValue string `xml:"Value" json:"currencyValue"`
}

type UpdateDatetime struct {
	Id             int    `sql:"id"`
	UpdateDatetime string `sql:"update_datetime"`
}
