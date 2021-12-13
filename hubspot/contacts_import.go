package hubspot

import (
	"strings"

	"github.com/grokify/goauth/scim"
	"github.com/grokify/gocharts/data/table"
	"github.com/grokify/mogo/strconv/phonenumber"
	"github.com/grokify/mogo/text/usstate"
	"github.com/nyaruka/phonenumbers"
)

func WriteContactsXLSX(filename string, users []scim.User) error {
	sheetdata := table.SheetData{
		SheetName: "Contacts",
		Rows: [][]interface{}{
			columnsInterface()}}

	for _, user := range users {
		row := userToScim(user)
		sheetdata.Rows = append(sheetdata.Rows, row)
	}
	return table.WriteXLSXInterface(filename, sheetdata)
}

func userToScim(user scim.User) []interface{} {
	row := []interface{}{
		user.Name.GivenName,
		user.Name.FamilyName,
		user.EmailAddress(),
		phonenumber.MustE164Format(user.PhoneNumber(), "US", phonenumbers.NATIONAL)}
	if len(user.Addresses) > 0 {
		addr := user.Addresses[0]
		row = append(row,
			addr.StreetAddress,
			addr.Locality,
			usstate.Abbreviate(addr.Region),
			addr.PostalCode)
	} else {
		row = append(row, "", "", "", "")
	}
	return row
}

const ColumnsString = `First Name,Last Name,Email Address,Phone Number,Street Address,City,State,Postal Code`

func Columns() []string {
	return strings.Split(ColumnsString, ",")
}

func columnsInterface() []interface{} {
	cols := []interface{}{}
	strs := Columns()
	for _, str := range strs {
		cols = append(cols, str)
	}
	return cols
}
