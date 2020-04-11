package entries

// Repository exposes a set of methods to store and retrieve data
type Repository interface {
	Insert(entry Entry) (bool, error)
	GetStateCumulative(state string) ([]Entry, error)
	GetCountryCumulative() ([]Entry, error)
	GetStateDeltas(state string) ([]Entry, error)
	GetCountryDeltas() ([]Entry, error)
	GetStateSummaryCumulative(state string) ([]Entry, error)
	GetCountrySummaryCumulative() ([]Entry, error)
	GetStateSummaryDeltas(state string) ([]Entry, error)
	GetCountrySummaryDeltas() ([]Entry, error)
	Close()
}
