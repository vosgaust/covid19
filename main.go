package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

type entry struct {
	date         string
	infected     int
	hospitalized int
	critical     int
	dead         int
	recovered    int
	active       int
}

func main() {
	csvFile, _ := os.Open("input.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	command := os.Args[1]
	option := ""
	if len(os.Args) > 2 {
		option = os.Args[2]
	}

	switch command {
	case "historic":
		if option != "" {
			getHistoricPerState(os.Args[2], aggregateByState(reader))
		} else {
			totals := aggreateByDate(reader)
			for _, value := range totals {
				fmt.Printf("%v %v %v %v %v\n", value.date, value.infected, value.dead, value.recovered, value.active)
			}
			// getHistoric(aggregateByState(reader))
		}
	case "summary":
		if option != "" {
			getSummaryPerState(os.Args[2], aggregateByState(reader))
		} else {
			for state := range aggregateByState(reader) {
				fmt.Printf("%v ", state)
				getSummaryPerState(state, aggregateByState(reader))
			}
		}
	}
}

func aggreateByDate(reader *csv.Reader) map[string]*entry {
	var totals = make(map[string]*entry)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if len(line) != 7 {
			continue
		} else if error != nil {
			log.Fatal(error)
		}

		date := line[1]
		infected := getNumber(line[2])
		dead := getNumber(line[5])
		recovered := getNumber(line[6])
		newEntry := entry{
			date:         line[1],
			infected:     infected,
			hospitalized: getNumber(line[3]),
			critical:     getNumber(line[4]),
			dead:         dead,
			recovered:    recovered,
			active:       infected - dead - recovered}
		if totals[date] != nil {
			totals[date].infected += newEntry.infected
			totals[date].hospitalized += newEntry.hospitalized
			totals[date].critical += newEntry.critical
			totals[date].dead += newEntry.dead
			totals[date].recovered += newEntry.recovered
			totals[date].active += newEntry.active
		} else {
			totals[date] = &newEntry
		}
	}
	return totals
}

func aggregateByState(reader *csv.Reader) map[string][]entry {
	var totals = make(map[string][]entry)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if len(line) != 7 {
			continue
		} else if error != nil {
			log.Fatal(error)
		}

		state := stateCodes[line[0]]
		infected := getNumber(line[2])
		dead := getNumber(line[5])
		recovered := getNumber(line[6])
		newEntry := entry{
			date:         line[1],
			infected:     infected,
			hospitalized: getNumber(line[3]),
			critical:     getNumber(line[4]),
			dead:         dead,
			recovered:    recovered,
			active:       infected - dead - recovered}
		if totals[state] != nil {
			totals[state] = append(totals[state], newEntry)
		} else {
			totals[state] = []entry{newEntry}
		}
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

func getHistoricPerState(state string, totals map[string][]entry) {
	for _, day := range totals[state] {
		fmt.Printf("%v %v %v %v %v\n", day.date, day.infected, day.dead, day.recovered, day.active)
	}
}

func getHistoric(totals map[string][]entry) {
	fmt.Println()
}

func getSummaryPerState(state string, totals map[string][]entry) {
	days := len(totals[state])
	today := totals[state][days-1]
	yesterday := totals[state][days-2]
	delta := entry{date: today.date,
		infected:     today.infected - yesterday.infected,
		hospitalized: today.hospitalized - yesterday.hospitalized,
		critical:     today.critical - yesterday.critical,
		dead:         today.dead - yesterday.dead,
		recovered:    today.recovered - yesterday.recovered,
		active:       today.active - yesterday.active}
	fmt.Printf("%v %v %v %v %v\n", delta.date, delta.infected, delta.dead, delta.recovered, delta.active)
}
