package entity

// MangaStatsHistory is entity for manga stats history.
type MangaStatsHistory struct {
	MangaID    int64
	Mean       float64
	Rank       int
	Popularity int
	Member     int
	Voter      int
}
