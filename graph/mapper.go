package graph

import (
	"fmt"
	"strconv"

	"github.com/ne241099/daifugo-server/graph/model"
	domain "github.com/ne241099/daifugo-server/model"
)

// parseID は GraphQL の文字列 ID を int64 に変換する
func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id %q: %w", s, err)
	}
	return id, nil
}

// mapUserToGraphQL はドメインの User を GraphQL モデルへ変換する
func mapUserToGraphQL(u *domain.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:        strconv.FormatInt(u.ID, 10),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// mapRoomToGraphQL はドメインの Room を GraphQL モデルへ変換する
func mapRoomToGraphQL(r *domain.Room) *model.Room {
	gRoom := &model.Room{
		ID:        strconv.FormatInt(r.ID, 10),
		Name:      r.Name,
		OwnerID:   strconv.FormatInt(r.OwnerID, 10),
		MemberIDs: make([]string, len(r.MemberIDs)),
		Game:      r.Game,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	for i, mid := range r.MemberIDs {
		gRoom.MemberIDs[i] = strconv.FormatInt(mid, 10)
	}
	return gRoom
}
