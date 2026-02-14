package tracker

type Player struct {
	Nickname       string             `json:"player"`
	Points         PointsInfo         `json:"points"`
	FavoriteServer FavoriteServerInfo `json:"favorite_server"`
}

type MapInfo struct {
	Points        int     `json:"points"`
	TotalFinishes int     `json:"total_finishes"`
	Finishes      int     `json:"finishes"`
	Rank          int     `json:"rank"`
	Time          float64 `json:"time"`
	FirstFinish   int64   `json:"first_finish"`
}

type Category struct {
	Points int                `json:"points"`
	Maps   map[string]MapInfo `json:"maps"`
}

type PointsInfo struct {
	Rank        int `json:"rank"`
	Points      int `json:"points"`
	TotalPoints int `json:"total"`
}

type FavoriteServerInfo struct {
	FavoriteServer string `json:"server"`
}
