package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/vosgaust/covid19/entries"
)

var stateCodes = map[string]string{
	"AN": "Andalucia",
	"AR": "Aragon",
	"AS": "Asturias",
	"CN": "Canarias",
	"CB": "Cantabria",
	"CM": "Castilla-Lamancha",
	"CL": "Castilla-Leon",
	"CT": "Catalunya",
	"EX": "Extremadura",
	"GA": "Galicia",
	"IB": "Baleares",
	"RI": "Rioja",
	"MD": "Madrid",
	"MC": "Murcia",
	"NC": "Navarra",
	"PV": "PaisVasco",
	"VC": "Valenciana",
	"ML": "Melilla",
	"CE": "Ceuta"}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	entriesRepository := entries.NewMySQL(cfg.MySQL)
	entriesService := entries.NewService(entriesRepository)
	defer entriesRepository.Close()

	if len(os.Args) < 2 {
		run(cfg)
	} else {
		command := os.Args[1]
		subcommand := ""
		if len(os.Args) > 2 {
			subcommand = os.Args[2]
		}
		state := ""
		if len(os.Args) > 3 {
			state = os.Args[3]
		}

		csvFile, _ := os.Open("input.csv")
		reader := csv.NewReader(bufio.NewReader(csvFile))

		switch command {
		case "historic":
			switch subcommand {
			case "deltas":
				result, err := entriesService.GetDeltas(state)
				if err != nil {
					panic(err.Error())
				}
				printResult(result)
			case "cumulative":
				result, err := entriesService.GetCumulative(state)
				if err != nil {
					panic(err.Error())
				}
				printResult(result)
			}
		case "summary":
			switch subcommand {
			case "deltas":
				if state == "" {
					fmt.Println(entriesRepository.GetCountrySummaryDeltas())
				} else {
					fmt.Println(entriesRepository.GetStateSummaryDeltas(state))
				}
			case "cumulative":
				if state == "" {
					fmt.Println(entriesRepository.GetCountrySummaryCumulative())
				} else {
					fmt.Println(entriesRepository.GetStateSummaryCumulative(state))
				}
			}
		case "collect":
			entries := parseCSV(reader)
			totalDeltas := processDeltas(entries)
			for _, entry := range totalDeltas {
				_, err = entriesRepository.Insert(entry)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func printResult(entries []entries.Entry) {
	for _, entry := range entries {
		fmt.Printf("%s %s %d %d %d %d\n", entry.Date, entry.Country, entry.Infected, entry.Dead, entry.Recovered, entry.Active)
	}
}

func processDeltas(totals []entries.Entry) []entries.Entry {
	var lastAdded = make(map[string]entries.Entry)
	var result []entries.Entry
	for _, entry := range totals {
		resultingEntry := entry
		if lastEntry, ok := lastAdded[entry.State]; ok {
			resultingEntry = entries.Entry{
				Date:         entry.Date,
				State:        entry.State,
				Country:      entry.Country,
				Infected:     entry.Infected - lastEntry.Infected,
				Hospitalized: entry.Hospitalized - lastEntry.Hospitalized,
				Critical:     entry.Critical - lastEntry.Critical,
				Dead:         entry.Dead - lastEntry.Dead,
				Recovered:    entry.Recovered - lastEntry.Recovered,
				Active:       entry.Active - lastEntry.Active}
		}
		fmt.Println(resultingEntry)
		lastAdded[entry.State] = entry
		result = append(result, resultingEntry)
	}
	return result
}

func formatDate(date string) string {
	layout := "2/1/2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println(err)
	}
	return t.Format("2006-01-02")
}

func parseCSV(reader *csv.Reader) []entries.Entry {
	var totals []entries.Entry
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if len(line) != 7 {
			continue
		} else if error != nil {
			log.Fatal(error)
		}

		date := formatDate(line[1])
		infected := getNumber(line[2])
		dead := getNumber(line[5])
		recovered := getNumber(line[6])
		newEntry := entries.Entry{
			Date:         date,
			Country:      "SP",
			State:        line[0],
			Infected:     infected,
			Hospitalized: getNumber(line[3]),
			Critical:     getNumber(line[4]),
			Dead:         dead,
			Recovered:    recovered,
			Active:       infected - dead - recovered}
		totals = append(totals, newEntry)
	}
	return totals
}

func getNumber(input string) int {
	if input == "" {
		return 0
	}
	number, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return number
}
