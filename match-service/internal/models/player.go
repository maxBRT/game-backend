package models

type Player struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"Role"`
	TicketID string `json:"ticketID"`
}

type PlayerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
