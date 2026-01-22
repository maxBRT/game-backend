package models

type Match struct {
	ID        string   `json:"id"`
	Survivors []Player `json:"survivors"`
	Killer    Player   `json:"killer"`
	CreatedAt string   `json:"createdAt"`
}
