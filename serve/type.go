package serve

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationData int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relatoin struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Locations struct {
	Id       int      `json:"id"`
	Location []string `json:"locations"`
}

type Dates struct {
	Id   int      `json:"id"`
	Date []string `json:"dates"`
}

type Artis struct {
	Artist   Artist
	Location Locations
	Date     Dates
	Relatoin Relatoin
}

type Serve struct {
	Port string
}

func NewPort(port string) *Serve {
	return &Serve{
		Port: port,
	}
}
