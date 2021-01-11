package main

import genshinapi "github.com/darylhjd/genshin-api"

// PageData : All page data structs will include this
type PageData struct {
	Title string
}

// IndexPageData : Struct containing required data for index page.
type IndexPageData struct {
	PageData
	DataTypes []string
}

// DataListPageData : Struct containing required data for data type page.
type DataListPageData struct {
	PageData
	DataType string
	NameList []string
}

// DataEntryPageData : Struct containing required data for a data entry page.
type DataEntryPageData struct {
	PageData
	Data genshinapi.DataEntry
}
