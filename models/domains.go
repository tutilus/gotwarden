package models

// Domains don't what is it
type Domains struct {
	EquivalentDomains       []byte `json:"EquivalentDomains"`
	GlobalEquivalentDomains []struct {
		Type     int      `json:"Type"`
		Domains  []string `json:"Domains"`
		Excluded bool     `json:"Excluded"`
	} `json:"GlobalEquivalentDomains"`
	Object string `json:"Object"`
}
