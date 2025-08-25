package domain

type GachaItemWeights []struct {
	ItemId int64
	Weight int
}

type Gacha struct {
	Weights GachaItemWeights
}
