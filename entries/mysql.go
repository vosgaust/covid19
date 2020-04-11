package entries

import (
	"database/sql"
	"fmt"
	"log"
)

type mySQLRepository struct {
	client *sql.DB
	table  string
}

// NewMySQL creates a new instance of a mysql repository
func NewMySQL(config Config) Repository {
	db, err := sql.Open("mysql", config.Connection)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	return &mySQLRepository{db, config.Table}
}

func (repository *mySQLRepository) Insert(entry Entry) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (date, country, state, infected, hospitalized, critical, dead, recovered, active) VALUES(\"%s\", \"%s\", \"%s\", %d, %d, %d, %d, %d, %d)", repository.table, entry.Date, entry.Country, entry.State, entry.Infected, entry.Hospitalized, entry.Critical, entry.Dead, entry.Recovered, entry.Active)
	insert, err := repository.client.Query(query)

	if err != nil {
		return false, err
	}
	defer insert.Close()

	return true, nil
}

func (repository *mySQLRepository) Close() {
	repository.client.Close()
}

func (repository *mySQLRepository) GetStateCumulative(state string) ([]Entry, error) {
	var result []Entry
	stateQuery := fmt.Sprintf("select * from %s where state='%s'", repository.table, state)
	query := fmt.Sprintf("SELECT a.date, a.country, a.state, sum(b.infected), sum(b.hospitalized), sum(b.critical), sum(b.dead), sum(b.recovered), sum(b.active) AS cumulative FROM (%s) a JOIN (%s) b ON a.date >= b.date GROUP BY date,country;", stateQuery, stateQuery)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.State, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetCountryCumulative() ([]Entry, error) {
	var result []Entry
	countryQuery := fmt.Sprintf("select date, country, sum(infected) as infected, sum(hospitalized) as hospitalized, sum(critical) as critical, sum(dead) as dead, sum(recovered) as recovered, sum(active) as active from %s group by date,country", repository.table)
	query := fmt.Sprintf("SELECT a.date, a.country, sum(b.infected), sum(b.hospitalized), sum(b.critical), sum(b.dead), sum(b.recovered), sum(b.active) AS cumulative FROM (%s) a JOIN (%s) b ON a.date >= b.date GROUP BY date,country;", countryQuery, countryQuery)
	fmt.Println(query)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetCountryDeltas() ([]Entry, error) {
	var result []Entry
	query := fmt.Sprintf("select date, country, sum(infected), sum(hospitalized), sum(critical), sum(dead), sum(recovered), sum(active) from %s group by date,country;", repository.table)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		fmt.Printf("failing: %v", err.Error())
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetStateDeltas(state string) ([]Entry, error) {
	var result []Entry
	query := fmt.Sprintf("select date, country, state, infected, hospitalized, critical, dead, recovered, active from %s where state='%s';", repository.table, state)
	queryResult, err := repository.client.Query(query)
	fmt.Println(query)
	if err != nil {
		fmt.Printf("failing: %v", err.Error())
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.State, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetStateSummaryCumulative(state string) ([]Entry, error) {
	var result []Entry
	stateQuery := fmt.Sprintf("select * from %s where state='%s'", repository.table, state)
	query := fmt.Sprintf("SELECT * FROM (SELECT a.date, a.country, a.state, sum(b.infected), sum(b.hospitalized), sum(b.critical), sum(b.dead), sum(b.recovered), sum(b.active) AS cumulative FROM (%s) a JOIN (%s) b ON a.date >= b.date GROUP BY date,country) AS total WHERE date='2020-04-10';", stateQuery, stateQuery)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.State, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetCountrySummaryCumulative() ([]Entry, error) {
	var result []Entry
	countryQuery := fmt.Sprintf("select date, country, sum(infected) as infected, sum(hospitalized) as hospitalized, sum(critical) as critical, sum(dead) as dead, sum(recovered) as recovered, sum(active) as active from %s group by date,country", repository.table)
	query := fmt.Sprintf("SELECT * FROM (SELECT a.date, a.country, sum(b.infected), sum(b.hospitalized), sum(b.critical), sum(b.dead), sum(b.recovered), sum(b.active) AS cumulative FROM (%s) a JOIN (%s) b ON a.date >= b.date GROUP BY date,country) AS total WHERE date='2020-04-10';", countryQuery, countryQuery)
	fmt.Println(query)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetStateSummaryDeltas(state string) ([]Entry, error) {
	var result []Entry
	query := fmt.Sprintf("SELECT date, country, state, infected, hospitalized, critical, dead, recovered, active FROM %s WHERE state='%s' AND date='2020-04-10';", repository.table, state)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.State, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}

func (repository *mySQLRepository) GetCountrySummaryDeltas() ([]Entry, error) {
	var result []Entry
	query := fmt.Sprintf("select date, country, sum(infected), sum(hospitalized), sum(critical), sum(dead), sum(recovered), sum(active) from %s WHERE date='2020-04-10' GROUP BY country;", repository.table)
	queryResult, err := repository.client.Query(query)
	if err != nil {
		return nil, err
	}
	for queryResult.Next() {
		var newEntry Entry
		err = queryResult.Scan(&newEntry.Date, &newEntry.Country, &newEntry.Infected, &newEntry.Hospitalized, &newEntry.Critical, &newEntry.Dead, &newEntry.Recovered, &newEntry.Active)
		if err != nil {
			log.Println(err.Error())
		}
		result = append(result, newEntry)
	}

	defer queryResult.Close()

	return result, nil
}
