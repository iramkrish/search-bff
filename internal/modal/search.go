package model

type Hotel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Flight struct {
	ID     string `json:"id"`
	Number string `json:"number"`
}

type SearchResponse struct {
	Hotels  []Hotel  `json:"hotels,omitempty"`
	Flights []Flight `json:"flights,omitempty"`
}
