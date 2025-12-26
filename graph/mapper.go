package graph

import (
	"strconv"

	"github.com/ne241099/daifugo-server/graph/model"
	"github.com/ne241099/daifugo-server/internal/game"
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
		gRoom.Game = &model.Game{
			IsRevolution: r.Game.IsRevolution,
			FieldCards:   mapCards(r.Game.FieldCards),
			Players:      make([]*model.GamePlayer, len(r.Game.Players)),
		}

		turnPlayer := r.Game.Players[r.Game.Turn]
		gRoom.Game.TurnUserID = strconv.FormatInt(turnPlayer.UserID, 10)

		for i, p := range r.Game.Players {
			gRoom.Game.Players[i] = &model.GamePlayer{
				UserID: strconv.FormatInt(p.UserID, 10),
				Hand:   mapCards(p.Hand),
				Rank:   int32(p.Rank),
			}
		}
	}

	return gRoom
}

// ドメインのCardモデルリストをGraphQLのCardモデルリストに変換
func mapCards(cards []*game.Card) []*model.Card {
	res := make([]*model.Card, len(cards))
	for i, c := range cards {
		res[i] = &model.Card{
			ID:   int32(c.ID),
			Suit: c.Suit.String(),
			Rank: int32(c.Rank),
		}
	}
	return res
}
