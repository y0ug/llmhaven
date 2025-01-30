package examples

import (
	"github.com/invopop/jsonschema"
	"github.com/y0ug/llmhaven/chat"
)

type GetCoordinatesInput struct {
	Location string `json:"location" jsonschema_description:"The location to look up."`
}

var GetCoordinatesInputSchema = GenerateSchema[GetCoordinatesInput]()

type GetCoordinateResponse struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

func GetCoordinates(location string) GetCoordinateResponse {
	return GetCoordinateResponse{
		Long: -122.4194,
		Lat:  37.7749,
	}
}

// Get Temperature Unit

type GetTemperatureUnitInput struct {
	Country string `json:"country" jsonschema_description:"The country"`
}

var GetTemperatureUnitInputSchema = GenerateSchema[GetTemperatureUnitInput]()

func GetTemperatureUnit(country string) string {
	return "farenheit"
}

// Get Weather

type GetWeatherInput struct {
	Lat  float64 `json:"lat"  jsonschema_description:"The latitude of the location to check weather."`
	Long float64 `json:"long" jsonschema_description:"The longitude of the location to check weather."`
	Unit string  `json:"unit" jsonschema_description:"Unit for the output"`
}

var GetWeatherInputSchema = GenerateSchema[GetWeatherInput]()

type GetWeatherResponse struct {
	Unit        string  `json:"unit"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(lat, long float64, unit string) GetWeatherResponse {
	return GetWeatherResponse{
		Unit:        "farenheit",
		Temperature: 122,
	}
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

func ToPtr[T any](s T) *T {
	return &s
}

func genTools() []chat.Tool {
	tools := []chat.Tool{
		{
			Name: "get_coordinates",
			Description: ToPtr(
				"Accepts a place as an address, then returns the latitude and longitude coordinates.",
			),
			InputSchema: GetCoordinatesInputSchema,
		},
		{
			Name:        "get_temperature_unit",
			InputSchema: GetTemperatureUnitInputSchema,
		},
		{
			Name:        "get_weather",
			Description: ToPtr("Get the weather at a specific location"),
			InputSchema: GetWeatherInputSchema,
		},
	}
	return tools
}
