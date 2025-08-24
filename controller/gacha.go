package controller

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"hex-arch-golang/db"
)

type Gacha struct {
	*sql.DB
}

func NewGach(db *sql.DB) *Gacha {
	return &Gacha{}
}

func (g *Gacha) Draw(ctx context.Context, gachaId int) (*db.Item, error) {
	gachaItems, err := db.GachaItems(
		qm.Select(db.GachaItemColumns.ItemID, db.GachaItemColumns.Weight),
		db.GachaItemWhere.GachaID.EQ(gachaId),
	).All(ctx, g.DB)

	if err != nil {
		return nil, err
	}

	weights := make([]int, len(gachaItems))
	for i, item := range gachaItems {
		weights[i] = item.Weight
	}

	seed := time.Now().UnixNano()

	index, err := linearSearchLottery(weights, seed)
	if err != nil {
		return nil, err
	}

	// Get Item Info
	item, err := db.FindItem(ctx, g.DB, gachaItems[index].ItemID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func linearSearchLottery(weights []int, seed int64) (int, error) {
	//  重みの総和を取得する
	var total int
	for _, weight := range weights {
		total += weight
	}

	// 乱数取得
	r := rand.New(rand.NewSource(seed))
	rnd := r.Intn(total)

	var currentWeight int
	for i, w := range weights {
		// 現在要素までの重みの総和
		currentWeight += w

		if rnd < currentWeight {
			return i, nil
		}
	}
	return 0, errors.New("the lottery failed")
}
