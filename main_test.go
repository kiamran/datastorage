package main

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
)

var streetResponse = `[
{
id: "STR005",
name: "Solomyanska",
prefix: "str.",
districts: [
{
id: "DST003",
name: null,
city: null
}
]
},
{
id: "STR002",
name: "Peremohy",
prefix: "ave",
districts: [
{
id: "DST003",
name: null,
city: null
}
]
}
]`

func TestDeserializingDistricts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://geo-data/streets/city/kyiv?locale=uk",
		httpmock.NewBytesResponder(200, []byte(streetResponse)))

	streets := getAllStreets("kyiv", "uk")

	for _, street := range streets {
		fmt.Println(street.Name)
	}
}
