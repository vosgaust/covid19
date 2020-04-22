package entries

// Service exposes methods to retrieve entries data
type Service interface {
	GetCumulative(state string) ([]Entry, error)
	GetDeltas(state string) ([]Entry, error)
}

type defaultService struct {
	repository Repository
}

// NewService creates a new instance of entries service
func NewService(repository Repository) Service {
	return &defaultService{repository}
}

func (s *defaultService) GetCumulative(state string) ([]Entry, error) {
	var result []Entry
	var err error
	if state == "" {
		result, err = s.repository.GetCountryCumulative()
	} else {
		result, err = s.repository.GetStateCumulative(state)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *defaultService) GetDeltas(state string) ([]Entry, error) {
	var result []Entry
	var err error
	if state == "" {
		result, err = s.repository.GetCountryDeltas()
	} else {
		result, err = s.repository.GetStateDeltas(state)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}
