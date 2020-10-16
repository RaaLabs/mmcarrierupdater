package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	origFileName := flag.String("fileName", "/usr/share/ModemManager/mm-dell-dw5821e-carrier-mapping.conf", "The filename of the carrier mapping file")
	newMappingsFile := flag.String("newMappingsFile", "mappings.txt", "Mappings in the format operatoridNR=apn, new mappings are separed by new lines")
	appendAfter := flag.String("appendAfter", "[dell dw5821e]", "The string in the original document to match, and start appending the new values after it")

	flag.Parse()

	mappingOrig, err := getOriginalMappings(*origFileName)
	if err != nil {
		fmt.Printf("error: getData failed: %v\n", err)
		return
	}

	mappingFromFile, err := getNewMappings(*newMappingsFile)
	if err != nil {
		fmt.Printf("error: getMappings failed: %v\n", err)
		return
	}

	missingMappings := findNotExistingMappings(mappingOrig, mappingFromFile)
	if missingMappings == nil {
		log.Printf("info: no new mappings, exiting...")
		return
	}

	err = updateMappings(*origFileName, mappingOrig, *appendAfter, missingMappings)
	if err != nil {
		log.Printf("error: updateMissingMappings failed: %v\n", err)
	}

}

func updateMappings(origFileName string, dataOrig []string, appendAfter string, missingMappings []string) error {
	var writePos int

	// Get the position of where to start writing after
	// If we find a match we set writePos to the index of
	// where in the slice the match was found, and if no
	// match was found we're setting the writePos to the
	// length of the slice so we do the appending after
	// the last element.
	for i, d := range dataOrig {
		if strings.Contains(d, appendAfter) {
			writePos = i
			fmt.Println(i)
		}
	}

	if writePos == 0 {
		log.Println("no match for the expression of where to append from, appending at the end")
		writePos = len(dataOrig)
	}

	var s []string
	s = append(s, dataOrig[:writePos+1]...)
	s = append(s, missingMappings...)
	// If no match was found for where to append from earlier we're defaulting to
	// append to the end of the slice. We then drop the last append so we don't
	// go over the bounds of the slice.
	if writePos != len(dataOrig) {
		s = append(s, dataOrig[writePos+1:]...)
	}

	fh, err := os.Create("test.txt")
	if err != nil {
		return err
	}
	defer fh.Close()

	for _, v := range s {
		_, err := fh.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// findNotExistingMappings compares the Original mapping file
// with the entries in the new mapping file, and returns the
// ones that are missing in the original file
func findNotExistingMappings(dataOrig []string, newMappings []string) []string {
	var newMappingsFound []string

	for _, m := range newMappings {
		var found bool

		for _, d := range dataOrig {
			if strings.Contains(d, m) {
				found = true
				break
			}
		}

		if found {
			continue
		}

		newMappingsFound = append(newMappingsFound, m)
	}

	return newMappingsFound
}

// getMappings will get the new mappings to check from file,
// and return a []string containing all the new mappings
func getNewMappings(fileName string) ([]string, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	var mappings []string

	for scanner.Scan() {
		s := scanner.Text()
		sp := strings.Split(s, " ")
		fmt.Println(sp)

		mappings = append(mappings, s)
	}

	return mappings, nil
}

// getData will read the original file content, split the the
// file line by line and put the content in a []string
func getOriginalMappings(fileName string) ([]string, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	var dataOrig []string

	for scanner.Scan() {
		dataOrig = append(dataOrig, scanner.Text())
	}

	return dataOrig, nil
}
