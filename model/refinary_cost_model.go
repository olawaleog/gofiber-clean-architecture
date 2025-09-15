package model

type RefineryCostModel struct {
	CostPerThousandLitre float64
	Distance             float64        `json:"distance"`
	TenThousandLitre     WaterCostModel `json:"tenThousandLitre"`
	TwentyThousandLitre  WaterCostModel `json:"twentyThousandLitre"`
	ThirtyThousandLitre  WaterCostModel `json:"thirtyThousandLitre"`
	FortyThousandLitre   WaterCostModel `json:"fortyThousandLitre"`
	Time                 int            `json:"time"`
	RefineryId           uint           `json:"refineryId"`
	Address              AddressModel   `json:"address"`
	Currency             string         `json:"currency"`
}
