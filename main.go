// Package familytreemap calculates relative population frequencies
// from Family Tree DNA projects.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/yogischogi/familytreemap/ftdna"
)

func main() {
	var (
		in       = flag.String("in", "", "Filename of input table.")
		out      = flag.String("out", "", "Filename for results in CSV format.")
		col      = flag.Int("col", 3, "Column number that contains country information.")
		relative = flag.Bool("relative", true, "Calculates relative or absolute frequencies.")
		tin      = flag.String("tin", "", "Testers in: A file that contains the number of persons who have tested for each country.")
		sumuk    = flag.Bool("sumuk", false, "Adds the number of testers from England, Wales, Scotland and Northern Ireland to United Kingdom.")
	)
	flag.Parse()

	// Get the number of testers for each country.
	var totalTesters map[string]float32
	var err error
	if *tin != "" {
		totalTesters, err = ftdna.ReadCountryTesters(*tin)
		if err != nil {
			fmt.Printf("Error reading the number of testers from file, %v.\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print("Please provide the number of total testers " +
			"for each country by using -tin <filename>.\r\n" +
			"Using example data instead.\r\n")
		totalTesters = ftdna.TotalTesters
	}

	// Read Family Tree project data.
	projectTable, err := ftdna.NewTableFromFile(*in, *col-1)
	if err != nil {
		fmt.Printf("Error reading project table, %v.\n", err)
		os.Exit(1)
	}

	countryFreqs := projectTable.FrequenciesOf(totalTesters)

	// Add testers from England, Wales, Scotland and N. Ireland to United Kingdom.
	if *sumuk == true {
		countryFreqs.SumUKTesters()
		if totalTesters != nil {
			totalTesters["United Kingdom"] += totalTesters["England"] +
				totalTesters["Wales"] +
				totalTesters["Scotland"] +
				totalTesters["Northern Ireland"]
		}
	}

	// Calculate relative frequencies.
	if *relative == true {
		for i, _ := range countryFreqs {
			countryFreqs[i].Persons = 100 * countryFreqs[i].Persons / totalTesters[countryFreqs[i].Country]
		}
	}

	sort.Stable(sort.Reverse(&countryFreqs))

	// Write results to file.
	err = countryFreqs.WriteCSV(*out)
	if err != nil {
		fmt.Printf("Error writing result to file in CSV format, %v.\r\n", err)
		os.Exit(1)
	}
}
