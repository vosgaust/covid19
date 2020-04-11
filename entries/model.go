package entries

// Entry defines the fields for a entry
type Entry struct {
	Date         string
	Country      string
	State        string
	Infected     int
	Hospitalized int
	Critical     int
	Dead         int
	Recovered    int
	Active       int
}
