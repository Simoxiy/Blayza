package serve

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
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
	err := GetData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	tmp, err := template.ParseFiles("./template/index.html")
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	if err := tmp.Execute(w, artists); err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
}

func Song(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	numbreid, err := strconv.Atoi(id)
	if len(r.URL.Query()) > 1 || numbreid > 52 || numbreid < 1 || err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}
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
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%v", id), &artis.Artist)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%v", id), &artis.Location)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%v", id), &artis.Date)
	fetcheData(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%v", id), &artis.Relatoin)
	wg.Wait()
	if artis.Artist.Id == 0 || artis.Date.Id == 0 || artis.Location.Id == 0 || artis.Relatoin.Id == 0 {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	tmp, err := template.ParseFiles("./template/artsit.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	if err := tmp.Execute(w, artis); err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}
	tmp, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	url := "https://groupietrackers.herokuapp.com/api/artists"
	rep, err := http.Get(url)
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	bytedata, err := io.ReadAll(rep.Body)
	if err != nil {
		http.Error(w, "internal server error ", http.StatusInternalServerError)
		return
	}
	search := fmt.Sprintf(`{[^{}]*"name":\s*"(?i)%v.*?"[^{}]*\}`, name)
	re := regexp.MustCompile(search)
	if !re.Match((bytedata)) {
		tmp.Execute(w, nil)
		return
	}
	searched := re.FindAll((bytedata), 52)
	artist := []Artist{}
	arts := Artist{}
	for i := 0; i < len(searched); i++ {
		err := json.Unmarshal((searched[i]), &arts)
		if err != nil {
			http.Error(w, "internal server error ", http.StatusInternalServerError)
			return
		}
		artist = append(artist, arts)
	}
	tmp.Execute(w, artist)
}

func Style(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/css/" {
		http.NotFound(w, r)
		return
	}
	styleserv := http.FileServer(http.Dir("style"))
	http.StripPrefix("/css", styleserv).ServeHTTP(w, r)
}

func GetData(url string, data any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("errer")
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(data)
}

func (r *Serve) Start() error {
	http.HandleFunc("/", Artists)
	http.HandleFunc("/art", Song)
	http.HandleFunc("/search", Search)
	http.HandleFunc("/api/", Api)
	http.HandleFunc("/api/{id}", GetArtist)
	http.HandleFunc("/css/", Style)
	return http.ListenAndServe(r.Port, nil)
}
