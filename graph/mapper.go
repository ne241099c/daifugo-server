package graph

import (
	"strconv"

	"github.com/ne241099/daifugo-server/graph/model"
	domain "github.com/ne241099/daifugo-server/model"
)

func mapRoomToGraphQL(r *domain.Room) *model.Room {
	gRoom := &model.Room{
		ID:        strconv.FormatInt(r.ID, 10),
		Name:      r.Name,
		OwnerID:   strconv.FormatInt(r.OwnerID, 10),
		MemberIds: make([]string, len(r.MemberIDs)),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	for i, mid := range r.MemberIDs {
		gRoom.MemberIds[i] = strconv.FormatInt(mid, 10)
	}

	if r.Game != nil {
		gRoom.Game = r.Game
	}

	return gRoom
}
