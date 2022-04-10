package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
)

func AthleteHandler(w http.ResponseWriter, r *http.Request) {
	var name string
	queries := r.URL.Query()

	if nameSlc, ok := queries["name"]; !ok || len(nameSlc[0]) <= 0 {
		http.Error(w, "no name param", http.StatusBadRequest)
		return
	} else {
		name = nameSlc[0]
	}

	filtered := Filter(athletes, func(athlete Athlete) bool {
		return athlete.Athlete == name
	})

	if len(filtered) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := athletToInfo(&filtered[0])

	for _, athlete := range filtered {
		_, ok := resp.MedalsByYear[athlete.Year]
		if !ok {
			resp.MedalsByYear[athlete.Year] = &Medals{0, 0, 0, 0}
		}
		resp.MedalsByYear[athlete.Year].Gold += athlete.Gold
		resp.MedalsByYear[athlete.Year].Silver += athlete.Silver
		resp.MedalsByYear[athlete.Year].Bronze += athlete.Bronze
		resp.MedalsByYear[athlete.Year].Total += athlete.Gold + athlete.Silver + athlete.Bronze

		resp.Medals.Gold += athlete.Gold
		resp.Medals.Silver += athlete.Silver
		resp.Medals.Bronze += athlete.Bronze
		resp.Medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func TopAthletesHandler(w http.ResponseWriter, r *http.Request) {
	var sportParam string
	var limitParam int
	var err error
	queries := r.URL.Query()

	if sportSlc, ok := queries["sport"]; ok {
		if sportSlc[0] != "" {
			sportParam = sportSlc[0]
		} else {
			http.Error(w, "no sport param", http.StatusBadRequest)
		}
	}

	if limitSlc, ok := queries["limit"]; ok {
		if limitSlc[0] != "" {
			limitParam, err = strconv.Atoi(limitSlc[0])
			if err != nil {
				http.Error(w, "invalid limit param", http.StatusBadRequest)
				return
			}
		}
	} else {
		limitParam = 3
	}

	test := func(athlete Athlete) bool {
		return athlete.Sport == sportParam
	}
	filtered := Filter(athletes, test)

	if len(filtered) == 0 {
		http.Error(w, "sport not found", http.StatusNotFound)
		return
	}

	filteredMap := GetAllAthlets(filtered)

	values := make([]*AthleteInfo, 0, len(filteredMap))

	for _, v := range filteredMap {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Medals.Gold != values[j].Medals.Gold {
			return values[i].Medals.Gold > values[j].Medals.Gold
		}
		if values[i].Medals.Silver != values[j].Medals.Silver {
			return values[i].Medals.Silver > values[j].Medals.Silver
		}
		if values[i].Medals.Bronze != values[j].Medals.Bronze {
			return values[i].Medals.Bronze > values[j].Medals.Bronze
		}

		return values[i].Athlete < values[j].Athlete
	})

	limit := int(math.Min(float64(limitParam), float64(len(values))))
	result := values[:limit]
	b, err := json.Marshal(&result)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func TopCountriesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("========")
	var yearParam int
	var limitParam int
	var err error
	queries := r.URL.Query()

	if yearSlc, ok := queries["year"]; ok {
		if yearSlc[0] != "" {
			yearParam, err = strconv.Atoi(yearSlc[0])
			if err != nil {
				http.Error(w, "invalid year param", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "no sport param", http.StatusBadRequest)
		}
	}

	if limitSlc, ok := queries["limit"]; ok {
		if limitSlc[0] != "" {
			limitParam, err = strconv.Atoi(limitSlc[0])
			if err != nil {
				http.Error(w, "invalid limit param", http.StatusBadRequest)
				return
			}
		}
	} else {
		limitParam = 3
	}

	test := func(athlete Athlete) bool {
		return athlete.Year == yearParam
	}
	filtered := Filter(athletes, test)

	if len(filtered) == 0 {
		http.Error(w, "year not found", http.StatusNotFound)
		return
	}

	filteredMap := GetAllCountries(filtered)

	values := make([]*CountryInfo, 0, len(filteredMap))

	for _, v := range filteredMap {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Gold != values[j].Gold {
			return values[i].Gold > values[j].Gold
		}
		if values[i].Silver != values[j].Silver {
			return values[i].Silver > values[j].Silver
		}
		if values[i].Bronze != values[j].Bronze {
			return values[i].Bronze > values[j].Bronze
		}

		return values[i].Country < values[j].Country
	})

	limit := int(math.Min(float64(limitParam), float64(len(values))))
	result := values[:limit]
	b, err := json.Marshal(&result)

	if err != nil {
		http.Error(w, "server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
