package models

const (
	Survivor = "survivor"
	Killer   = "killer"
)

func (p Player) IsSurvivor() bool {
	return p.Role == Survivor
}

func (p Player) IsKiller() bool {
	return p.Role == Killer
}
