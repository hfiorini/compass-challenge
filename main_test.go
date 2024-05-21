package main

import (
	"bytes"

	"os"
	"testing"

	"github.com/tealeg/xlsx"
)

func TestCalculateScoring(t *testing.T) {
	contact1 := Contact{ID: 1, Name: "John Doe", Name1: "JD", Email: "john@example.com", PostalZip: "12345", Address: "123 Elm St"}
	contact2 := Contact{ID: 2, Name: "John Doe", Name1: "JD", Email: "john@example.com", PostalZip: "12345", Address: "123 Elm St"}
	contact3 := Contact{ID: 3, Name: "Jane Doe", Name1: "JD", Email: "jane@example.com", PostalZip: "67890", Address: "456 Oak St"}

	tests := []struct {
		contact1 Contact
		contact2 Contact
		expected int
	}{
		{contact1, contact2, 5},
		{contact1, contact3, 1},
		{contact2, contact3, 1},
	}

	for _, test := range tests {
		score := calculateScoring(&test.contact1, &test.contact2)
		if score != test.expected {
			t.Errorf("Expected scoring %d but got %d", test.expected, score)
		}
	}
}

func TestFindDuplicates(t *testing.T) {
	contacts := []Contact{
		{ID: 1, Name: "John Doe", Name1: "JD", Email: "john@example.com", PostalZip: "12345", Address: "123 Elm St"},
		{ID: 2, Name: "Pepe Doe", Name1: "JD", Email: "john@example.com", PostalZip: "12345", Address: "123 Elm St"},
		{ID: 3, Name: "Jane Doe", Name1: "JD", Email: "jane@example.com", PostalZip: "67890", Address: "456 Oak St"},
	}

	expected := map[int][]Row{
		1: {
			{SourceID: 1, MatchID: 2, Scoring: 4},
			{SourceID: 1, MatchID: 3, Scoring: 1},
		},
		2: {
			{SourceID: 2, MatchID: 3, Scoring: 1},
		},
	}

	result := findDuplicates(contacts)

	if len(result) != len(expected) {
		t.Errorf("Expected %d duplicates but got %d", len(expected), len(result))
	}

	for key, rows := range expected {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected key %d not found", key)
		}
		for i, row := range rows {
			if row != result[key][i] {
				t.Errorf("Expected row %v but got %v", row, result[key][i])
			}
		}
	}
}

func createTempXLSX(t *testing.T, contacts []Contact) string {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		t.Fatal(err)
	}

	headerRow := sheet.AddRow()
	headers := []string{"ID", "Name", "Name1", "Email", "PostalZip", "Address"}
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	for _, contact := range contacts {
		row := sheet.AddRow()
		row.AddCell().SetInt(contact.ID)
		row.AddCell().SetString(contact.Name)
		row.AddCell().SetString(contact.Name1)
		row.AddCell().SetString(contact.Email)
		row.AddCell().SetString(contact.PostalZip)
		row.AddCell().SetString(contact.Address)
	}

	buf := new(bytes.Buffer)
	err = file.Write(buf)
	if err != nil {
		t.Fatal(err)
	}

	filename := "test_contacts.xlsx"
	err = file.Save(filename)
	if err != nil {
		t.Fatal(err)
	}

	return filename
}

func TestReadContactsFromXLSX(t *testing.T) {
	contacts := []Contact{
		{ID: 1, Name: "John Doe", Name1: "JD", Email: "john@example.com", PostalZip: "12345", Address: "123 Elm St"},
		{ID: 2, Name: "Jane Doe", Name1: "JD", Email: "jane@example.com", PostalZip: "67890", Address: "456 Oak St"},
	}

	filename := createTempXLSX(t, contacts)

	readContacts, err := readContactsFromXLSX(filename)
	if err != nil {
		t.Fatal(err)
	}

	if len(readContacts) != len(contacts) {
		t.Errorf("Expected %d contacts but got %d", len(contacts), len(readContacts))
	}

	for i, contact := range contacts {
		if readContacts[i] != contact {
			t.Errorf("Expected contact %v but got %v", contact, readContacts[i])
		}
	}
}

func TestPrintResults(t *testing.T) {
	duplicates := map[int][]Row{
		1: {
			{SourceID: 1, MatchID: 2, Scoring: 5, Accuracy: "High"},
			{SourceID: 1, MatchID: 3, Scoring: 1, Accuracy: "Low"},
		},
	}

	expectedOutput := "| Source ID |  Match | Accuracy | Scoring # \n" +
		"| 1 | 2 | High | 5\n" +
		"| 1 | 3 | Low | 1\n"

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdOut := os.Stdout
	os.Stdout = w

	printResults(duplicates)

	w.Close()
	os.Stdout = stdOut

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatal(err)
	}

	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected output:\n%s\nbut got:\n%s", expectedOutput, actualOutput)
	}
}
