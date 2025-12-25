package room

import (
	"context"

	"github.com/ne241099/daifugo-server/graph/model"
)

type CreateRoomUsecase interface {
	Execute(ctx context.Context, name string, ownerID int64, memberIDs []int64) (*model.Room, error)
}
