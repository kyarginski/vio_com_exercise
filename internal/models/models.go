package models

// Location represents location.
// swagger:model
type Location struct {
	IPAddress    string  `json:"ip_address"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue int64   `json:"mystery_value"`
}

// LoadStatistics information about load data.
type LoadStatistics struct {
	LoadTime   string `json:"load_time"`
	FilesCount int64  `json:"files_count"`
	Accepted   int64  `json:"accepted"`
	Discarded  int64  `json:"discarded"`
	Total      int64  `json:"total"`
}
