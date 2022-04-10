package main

type Athlete struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type Medals struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type AthleteInfo struct {
	Athlete      string          `json:"athlete"`
	Country      string          `json:"country"`
	Medals       Medals          `json:"medals"`
	MedalsByYear map[int]*Medals `json:"medals_by_year"`
}

type CountryInfo struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

func athletToInfo(a *Athlete) AthleteInfo {
	var info AthleteInfo
	info.Athlete = a.Athlete
	info.MedalsByYear = make(map[int]*Medals)
	info.Country = a.Country

	return info
}
