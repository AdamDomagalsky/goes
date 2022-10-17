package entities

// The Book struct defines the data associated with a
// book that the library own.
type Book struct {
	ID            int
	Title         string
	Author        string
	ISBN          string
	YearPublished int
	Description   string
	IsBorrowable  bool
}
