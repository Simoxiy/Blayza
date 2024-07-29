package serve

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func Artists(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "mothod not allow", http.StatusMethodNotAllowed)
		return
	}
	artists := []Artist{}
	GetData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	// w.Header().Add("Content-type", "application/json")
	// json.NewEncoder(w).Encode(artists)
	tmp, err := template.ParseFiles("./template/index.html")
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, artists)
}

func Song(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusAccepted)
	}
	// var wg sync.WaitGroup
	artis := &Artis{}

	fetcheData := func(url string, data any) {
		// defer wg.Done()
		err := GetData(url, data)
		if err != nil {
			return
		}
	}
	// var err error
	// wg.Add(4)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%v", id), &artis.Artist)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%v", id), &artis.Location)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%v", id), &artis.Date)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%v", id), &artis.Relatoin)
	// wg.Wait()
	if artis.Artist.Id == 0 || artis.Date.Id == 0 || artis.Location.Id == 0 || artis.Relatoin.Id == 0 {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(artis)
}

func GetData(url string, data any) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	} else if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return fmt.Errorf("errer")
	}
	defer resp.Body.Close()
	fmt.Println(json.NewDecoder(resp.Body).Decode(data))
	return nil
}

func (r *Serve) Start() error {
	http.HandleFunc("/", Artists)
	http.HandleFunc("/art", Song)
	return http.ListenAndServe(r.Port, nil)
}
