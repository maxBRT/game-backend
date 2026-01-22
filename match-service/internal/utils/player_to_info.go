package utils

import (
	m "github.com/maxbrt/game-backend/match-service/internal/models"
)

func PlayerToInfo(player m.Player) m.PlayerInfo {
	return m.PlayerInfo{
		ID:   player.ID,
		Name: player.Name,
		Role: player.Role,
	}
}

func PlayersToInfo(players []m.Player) []m.PlayerInfo {
	var infos []m.PlayerInfo
	for _, player := range players {
		infos = append(infos, PlayerToInfo(player))
	}
	return infos
}
