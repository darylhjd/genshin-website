package main

import (
	"encoding/json"
	"fmt"
	genshinapi "github.com/darylhjd/genshin-api"
	"log"
	"net/http"
	"strings"
)

// IndexHandler : Handler for index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	log.Println("IndexHandler: ", paths)
	// Since mux matches this function for any other
	// patterns that are not registered, we need to check
	// if the path is the exact match for the index page.
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}

	dataTypes, err := genshinapi.GetDataTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := IndexPageData{
		PageData:  PageData{Title: "genshin.dev"},
		DataTypes: dataTypes,
	}
	w.WriteHeader(http.StatusOK)
	err = indexTemplate.ExecuteTemplate(w, "index.html", &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NotFoundHandler : Handler for when a page does not exist.
func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("Sorry, this page does not exist!"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// MakeDataTypeHandler : Function closure to create handler for each data type.
func MakeDataTypeHandler(dtype string) (string, func(http.ResponseWriter, *http.Request)) {
	return fmt.Sprintf("/%s/", dtype), func(w http.ResponseWriter, r *http.Request) {
		paths := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		log.Printf("%sHandler: %v\n", strings.Title(dtype), paths)
		if len(paths) > 1 {
			GetSpecificDataEntryPage(w, paths)
		} else {
			GetDataListPage(w, dtype)
		}
	}
}

// GetSpecificDataEntryPage : Response to request to getting a specific data entry.
func GetSpecificDataEntryPage(w http.ResponseWriter, datas []string) {
	log.Printf("Getting specific %s '%s'\n", datas[0], datas[1])
	var dataentry genshinapi.DataEntry
	switch datas[0] {
	case genshinapi.ArtifactsDType:
		dataentry = new(genshinapi.Artifact)
	case genshinapi.CharactersDType:
		dataentry = new(genshinapi.Character)
	case genshinapi.DomainsDType:
		dataentry = new(genshinapi.Domain)
	case genshinapi.ElementsDType:
		dataentry = new(genshinapi.Element)
	case genshinapi.NationsDType:
		dataentry = new(genshinapi.Nation)
	case genshinapi.WeaponsDType:
		dataentry = new(genshinapi.Weapon)
	default:
		http.Error(w, "Unknown data type requested!", http.StatusInternalServerError)
		return
	}

	bytedata, err := genshinapi.GetCustomBody(datas...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytedata, dataentry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DataEntryPageData{
		PageData: PageData{Title: datas[1]},
		Data:     dataentry,
	}
	w.WriteHeader(http.StatusOK)
	err = dataEntryTemplate.ExecuteTemplate(w,
		fmt.Sprintf("%s.html", strings.TrimSuffix(datas[0], "s")), &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetDataListPage : Response to request to getting a list of entries for each data type.
func GetDataListPage(w http.ResponseWriter, dtype string) {
	var names []string
	bytedata, err := genshinapi.GetCustomBody(dtype)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytedata, &names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DataListPageData{
		PageData: PageData{Title: strings.Title(dtype)},
		DataType: dtype,
		NameList: names,
	}
	w.WriteHeader(http.StatusOK)
	err = dataTypeTemplate.ExecuteTemplate(w, fmt.Sprintf("%s.html", dtype), &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
