package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"text/template"
	"strings"
)

type District struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

type Street struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	DisplayName string `json:"displayName"`
	Districts   []District `json:"Districts,omitempty"`
}

type DatastorageFile struct {
	Districts []District `json:"districts"`
	Streets   []Street `json:"streets"`
}

var geoDataUrl = "http://geo-data"

func handler(w http.ResponseWriter, r *http.Request) {
	cityId := r.URL.Query().Get("cityId")
	locale := r.URL.Query().Get("locale")

	files, err := template.ParseFiles("./datastorage.template")

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=datastorage.js")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	marshal, err := json.Marshal(getDataStorage(cityId, locale))

	if err == nil {
		files.Execute(w, string(marshal[:]))
	}

}
func getDataStorage(cityId string, locale string) DatastorageFile {
	d := DatastorageFile{Streets: getAllStreets(cityId, locale), Districts: getAllDistricts(cityId, locale)}
	return d
}

func getAllDistricts(cityId, locale string) ([]District) {
	responses := make([]District, 0)
	//i := geoDataUrl + "/districts/city/" + cityId
	i := "http://mockbin.org/bin/a9b399cd-0b30-4c8d-8083-cb91a840efa0"
	json.Unmarshal(makeGetRequest(i), &responses)

	return responses
}

func getAllStreets(cityId, locale string) ([]Street) {
	responses := make([]Street, 0)
	i := geoDataUrl + "/streets/city/" + cityId + "?locale=" + locale
	json.Unmarshal(makeGetRequest(i), &responses)

	setDistricInfo(responses)
	return responses
}

func setDistricInfo(streets []Street) {
	m := make(map[string]string)
	for i := range streets {
		street := streets[i]
		_, exists := m[strings.ToLower(street.Name)]
		if !exists {
			street.Districts = nil
		}
	}
}

func makeGetRequest(url string) ([]byte) {
	log.Println("Executing get request with URL " + url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Unable to make GET request to " + url)
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Unable to prase json response")
		panic(err.Error())
	}

	// TODO  not sure if needed
	defer resp.Body.Close()

	return body
}

func main() {
	http.HandleFunc("/datastorage", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
