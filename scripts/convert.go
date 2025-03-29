package main
import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"log"
)
type KML struct {
	XMLName xml.Name `xml:"kml"`
	Placemarks []Placemark `xml:"Document>Placemark"`
}
type Placemark struct {
	Name string `xml:"name"`
	Geometry Geometry `xml:"Polygon>coordinates"`
}
type Geometry struct {
	Coordinates string `xml:",chardata"`
}
type GeoJSON struct {
	Type     string     `json:"type"`
	Features []Feature `json:"features"`
}
type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry  `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
func main() {
	kmlFile := "input.kml"
	data, err := ioutil.ReadFile(kmlFile)
	if err != nil {
		log.Fatalf("Error reading KML file: %v", err)
	}
	var kml KML
	err = xml.Unmarshal(data, &kml)
	if err != nil {
		log.Fatalf("Error unmarshalling KML: %v", err)
	}
	var geoJSON GeoJSON
	geoJSON.Type = "FeatureCollection"
	for _, placemark := range kml.Placemarks {
		var feature Feature
		feature.Type = "Feature"
		feature.Properties = map[string]interface{}{"name": placemark.Name}
		coords := placemark.Geometry.Coordinates
		var coordinates [][]float64
		for _, coord := range stringToCoords(coords) {
			coordinates = append(coordinates, coord)
		}
		feature.Geometry.Type = "Polygon"
		feature.Geometry.Coordinates = coordinates
		geoJSON.Features = append(geoJSON.Features, feature)
	}
	geoJSONFile := "output.geo"
	geoJSONData, err := json.MarshalIndent(geoJSON, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling GeoJSON: %v", err)
	}
	err = ioutil.WriteFile(geoJSONFile, geoJSONData, 0644)
	if err != nil {
		log.Fatalf("Error writing GeoJSON file: %v", err)
	}
	fmt.Println("Conversion complete! GeoJSON saved to", geoJSONFile)
}
func stringToCoords(coords string) [][]float64 {
	var result [][]float64
	for _, coord := range splitCoords(coords) {
		coordParts := splitCoord(coord)
		result = append(result, coordParts)
	}
	return result
}
func splitCoords(coords string) []string {
	return splitByDelimiter(coords, " ")
}
func splitCoord(coord string) []float64 {
	var result []float64
	parts := splitByDelimiter(coord, ",")
	for _, part := range parts {
		var value float64
		fmt.Sscanf(part, "%f", &value)
		result = append(result, value)
	}
	return result
}
func splitByDelimiter(str, delimiter string) []string {
	var result []string
	for _, part := range strings.Split(str, delimiter) {
		if len(part) > 0 {
			result = append(result, part)
		}
	}
	return result
}
