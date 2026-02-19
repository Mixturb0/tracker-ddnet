package tracker

import "time"

type PlayerStats struct {
	Player                 string              `json:"player"`
	Points                 Points              `json:"points"`
	TeamRank               Rank                `json:"team_rank"`
	Rank                   Rank                `json:"rank"`
	PointsLastYear         Points              `json:"points_last_year"`
	PointsLastMonth        Points              `json:"points_last_month"`
	PointsLastWeek         Points              `json:"points_last_week"`
	FavoriteServer         FavoriteServer      `json:"favorite_server"`
	FirstFinish            FirstFinish         `json:"first_finish"`
	LastFinishes           []Finish            `json:"last_finishes"`
	FavoritePartners       []Partner           `json:"favorite_partners"`
	Types                  map[string]GameType `json:"types"`
	Activity               []Activity          `json:"activity"`
	HoursPlayedPast365Days int                 `json:"hours_played_past_365_days"`
}

type PlayerGive struct {
	Player     string `json:"player"`
	Points     Points `json:"points"`
	Finish     int    `json:"finishes"`
	TotalHours int    `json:"hours"`
}

type LadderPlayer struct {
	Rank   int    `json:"rank"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type PlayerInfo struct {
	Time string `json:"time"`
	Name string `json:"name"`
}

// Структура самого сервера
type ServerInfo struct {
	ID      int          `json:"id"`
	Region  string       `json:"region"` // "EUR" или "RUS"
	Title   string       `json:"title"`
	Players []PlayerInfo `json:"players"`
}

type Points struct {
	Rank   *int `json:"rank"`
	Points int  `json:"points"`
	Total  int  `json:"total,omitempty"`
}

type Rank struct {
	Rank *int `json:"rank"`
}

type FavoriteServer struct {
	Server string `json:"server"`
}

type FirstFinish struct {
	Timestamp float64 `json:"timestamp"`
	Map       string  `json:"map"`
	Time      float64 `json:"time"`
}

type Finish struct {
	Timestamp float64 `json:"timestamp"`
	Map       string  `json:"map"`
	Time      float64 `json:"time"`
	Country   string  `json:"country"`
	Type      string  `json:"type"`
}

type Partner struct {
	Name     string `json:"name"`
	Finishes int    `json:"finishes"`
}

type GameType struct {
	Points   Points             `json:"points"`
	TeamRank Rank               `json:"team_rank"`
	Rank     Rank               `json:"rank"`
	Maps     map[string]MapInfo `json:"maps"`
}

type MapInfo struct {
	Points        int     `json:"points"`
	TotalFinishes int     `json:"total_finishes"`
	Finishes      int     `json:"finishes"`
	Rank          *int    `json:"rank,omitempty"`
	Time          float64 `json:"time,omitempty"`
	FirstFinish   float64 `json:"first_finish,omitempty"`
}

type Activity struct {
	Date        string `json:"date"`
	HoursPlayed int    `json:"hours_played"`
}

func (f *FirstFinish) GetTimestamp() time.Time {
	return time.Unix(int64(f.Timestamp), 0)
}

func (f *Finish) GetTimestamp() time.Time {
	return time.Unix(int64(f.Timestamp), 0)
}

func (m *MapInfo) GetFirstFinishTime() *time.Time {
	if m.FirstFinish == 0 {
		return nil
	}
	t := time.Unix(int64(m.FirstFinish), 0)
	return &t
}

func (a *Activity) GetDate() (time.Time, error) {
	return time.Parse("2006-02-01", a.Date)
}

func (p *PlayerStats) TotalHourse() int {
	TotalHours := 0
	for _, v := range p.Activity {
		TotalHours += v.HoursPlayed
	}
	return TotalHours
}

func (p *PlayerStats) TotalFinishes() int {
	TotalFinish := 0
	for _, v := range p.Types {
		for _, v2 := range v.Maps {
			TotalFinish += v2.Finishes
		}
	}
	return TotalFinish
}
