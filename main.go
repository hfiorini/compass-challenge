package main

import (
	"fmt"
	"sort"

	"github.com/tealeg/xlsx"
)

type Contact struct {
	ID        int
	Name      string
	Name1     string
	Email     string
	PostalZip string
	Address   string
}

type Row struct {
	SourceID int
	MatchID  int
	Scoring  int
	Accuracy string
}

// Function to parse the xlsx file and store the contacts in memory
func readContactsFromXLSX(filename string) ([]Contact, error) {
	var contacts []Contact

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {

			if rowIndex == 0 {
				continue
			}

			contact := Contact{}
			for cellIndex, cell := range row.Cells {
				switch cellIndex {
				case 0:
					contact.ID, _ = cell.Int()
				case 1:
					contact.Name = cell.String()
				case 2:
					contact.Name1 = cell.String()
				case 3:
					contact.Email = cell.String()
				case 4:
					contact.PostalZip = cell.String()
				case 5:
					contact.Address = cell.String()
				}
			}
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}

func main() {

	contacts, err := readContactsFromXLSX("contacts.xlsx")
	if err != nil {
		fmt.Println("Error reading contacts:", err)
		return
	}

	duplicates := findDuplicates(contacts)

	printResults(duplicates)
}

// I create a map, using the contact ID as key and all the matches as value
func findDuplicates(contacts []Contact) map[int][]Row {
	duplicates := make(map[int][]Row)

	for i := 0; i < len(contacts); i++ {
		for j := i + 1; j < len(contacts); j++ {
			if contacts[i].ID == -1 || contacts[j].ID == -1 {
				continue
			}
			scoring := calculateScoring(&contacts[i], &contacts[j])
			if scoring > 0 {
				value := Row{
					SourceID: contacts[i].ID,
					MatchID:  contacts[j].ID,
					Scoring:  scoring,
				}
				if _, ok := duplicates[contacts[i].ID]; ok {

					duplicates[contacts[i].ID] = append(duplicates[contacts[i].ID], value)
				} else {
					duplicates[contacts[i].ID] = []Row{value}
				}
			}
		}
	}

	return duplicates
}

// I decided to add one score point for each field that matches.
// When scoring is higher than 2, will be considered "High"
func calculateScoring(contact1, contact2 *Contact) int {
	score := 0

	if contact1.Name == contact2.Name {
		score++
	}
	if contact1.Name1 == contact2.Name1 {
		score++
	}
	if contact1.Email == contact2.Email {
		score++
	}
	if contact1.PostalZip == contact2.PostalZip {
		score++
	}
	if contact1.Address == contact2.Address && contact1.Address != "" {
		score++
	}

	return score
}

// This function shows the expected output sorted by highest accuracy (scoring)
func printResults(duplicates map[int][]Row) {

	fmt.Println("| Source ID |  Match | Accuracy | Scoring # ")
	rowList := []Row{}
	for _, rows := range duplicates {
		for _, row := range rows {

			accuracy := "Low"
			if row.Scoring > 2 {
				accuracy = "High"
			}
			row.Accuracy = accuracy
			rowList = append(rowList, row)

		}
	}

	sort.Slice(rowList, func(i, j int) bool {
		return rowList[i].Scoring > rowList[j].Scoring
	})
	for _, row := range rowList {
		fmt.Printf("| %d | %d | %s | %d\n", row.SourceID, row.MatchID, row.Accuracy, row.Scoring)
	}

}
