package main

import (
	"fmt"
	"os"
	"time"

	"github.com/grokify/goauth/hubspot"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/time/timeutil"
)

func main() {
	hclient := hubspot.NewClientAPIKey(os.Getenv("HUBSPOT_API_KEY"))
	dt := time.Now()
	err := hubspot.ContactsV3ExportWriteFiles(hclient,
		fmt.Sprintf("_hubspot_contacts_v3_%s_", dt.Format(timeutil.RFC3339FullDate)),
		&hubspot.ContactsListV3Opts{
			Limit:      100,
			Properties: []string{"firstname", "lastname", "email", "company", "linkedin_url"}})
	logutil.FatalErr(err)
	fmt.Println("DONE")
}
