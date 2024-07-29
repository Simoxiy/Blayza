package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/ysnbhb/group-tracker/serve"
)

func main() {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	rep, _ := http.Get(url)
	bytedata, _ := io.ReadAll(rep.Body)
	arg := os.Args[1]
	search := fmt.Sprintf(`{[^{}]*"name":\s*"(?i)%v.*?"[^{}]*\}`, arg)
	re := regexp.MustCompile(search)
	fmt.Println(re.Match((bytedata)))
	r := re.FindAll((bytedata), 52)
	fmt.Println(len(r))
	// fmt.Println(string(r))
	artist := []serve.Artist{}
	arts := serve.Artist{}
	for i := 0; i < len(r); i++ {
		fmt.Println(json.Unmarshal((r[i]), &arts))
		artist = append(artist, arts)
	}
	fmt.Println(artist)
}
