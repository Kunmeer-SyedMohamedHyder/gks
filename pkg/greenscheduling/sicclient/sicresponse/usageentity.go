package sicresponse

// UsageEntity represents a single entity in the usage response.
type UsageEntity struct {
	ID                         string   `json:"id"`
	Type                       string   `json:"type"`
	EntityID                   string   `json:"entityId"`
	EntityMake                 string   `json:"entityMake"`
	EntityModel                string   `json:"entityModel"`
	EntityType                 string   `json:"entityType"`
	EntitySerialNum            string   `json:"entitySerialNum"`
	EntityProductID            string   `json:"entityProductId"`
	EntityManufactureTimestamp string   `json:"entityManufactureTimestamp"`
	LocationName               *string  `json:"locationName"`
	LocationID                 *string  `json:"locationId"`
	LocationCity               *string  `json:"locationCity"`
	LocationState              *string  `json:"locationState"`
	LocationCountry            *string  `json:"locationCountry"`
	Name                       string   `json:"name"`
	CostUsd                    *float64 `json:"costUsd"`       // Nullable float64
	Co2eMetricTon              *float64 `json:"co2eMetricTon"` // Nullable float64
	Kwh                        *float64 `json:"kwh"`           // Nullable float64
}

// GetCostUsd returns the CostUsd value, or 0 if it is nil.
func (e *UsageEntity) GetCostUsd() float64 {
	if e.CostUsd == nil {
		return 0
	}
	return *e.CostUsd
}

// GetCo2eMetricTon returns the Co2eMetricTon value, or 0 if it is nil.
func (e *UsageEntity) GetCo2eMetricTon() float64 {
	if e.Co2eMetricTon == nil {
		return 0
	}
	return *e.Co2eMetricTon
}

// GetKwh returns the Kwh value, or 0 if it is nil.
func (e *UsageEntity) GetKwh() float64 {
	if e.Kwh == nil {
		return 0
	}
	return *e.Kwh
}

// UsageByEntityResponse represents the entire response from the usage API.
type UsageByEntityResponse struct {
	Items  []UsageEntity `json:"items"`
	Count  int           `json:"count"`
	Total  int           `json:"total"`
	Offset int           `json:"offset"`
}
