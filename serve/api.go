package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func Api(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
	}
	artsist := []Artist{}
	GetData("https://groupietrackers.herokuapp.com/api/artists", &artsist)
	json.NewEncoder(w).Encode(artsist)
}

func GetArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/")
	var wg sync.WaitGroup
	artis := &Artis{}
	fetcheData := func(url string, data any) {
		defer wg.Done()
		err := GetData(url, data)
		if err != nil {
			return
		}
	}
	wg.Add(4)
	go fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%v", id), &artis.Artist)
	go fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%v", id), &artis.Location)
	go fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%v", id), &artis.Date)
	go fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%v", id), &artis.Relatoin)
	wg.Wait()
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(artis)
}
