//go:build !solution

package shopfront

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type counters struct {
	rdb *redis.Client
}

func toStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func (c counters) GetItems(ctx context.Context, ids []ItemID, userID UserID) ([]Item, error) {
	pipe := c.rdb.Pipeline()

	viewCount := make([]*redis.IntCmd, len(ids))
	viewed := make([]*redis.BoolCmd, len(ids))
	for i, id := range ids {
		key := "item_" + toStr(int64(id))
		viewCount[i] = pipe.SCard(ctx, key)
		viewed[i] = pipe.SIsMember(ctx, key, toStr(int64(userID)))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]Item, len(ids))
	for i := range res {
		res[i].ViewCount = int(viewCount[i].Val())
		res[i].Viewed = viewed[i].Val()
	}

	return res, nil
}

func (c counters) RecordView(ctx context.Context, id ItemID, userID UserID) error {
	pipe := c.rdb.Pipeline()
	//pipe.SAdd(ctx, "user_"+strconv.FormatInt(int64(userID), 10), id)
	pipe.SAdd(ctx, "item_"+toStr(int64(id)), toStr(int64(userID)))
	_, err := pipe.Exec(ctx)
	return err
}

func New(rdb *redis.Client) Counters {
	return counters{rdb}
}
