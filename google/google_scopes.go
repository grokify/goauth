package google

import "google.golang.org/api/slides/v1"

const (
	ScopeDrive                 string = slides.DriveScope                 // See, edit, create, and delete all of your Google Drive files
	ScopeDriveFile             string = slides.DriveFileScope             // View and manage Google Drive files and folders that you have opened or created with this app
	ScopeDriveReadonly         string = slides.DriveReadonlyScope         // See and download all your Google Drive files
	ScopePresentations         string = slides.PresentationsScope         // View and manage your Google Slides presentations
	ScopePresentationsReadonly string = slides.PresentationsReadonlyScope // View your Google Slides presentations
	ScopeSpreadsheets          string = slides.SpreadsheetsScope          // See, edit, create, and delete your spreadsheets in Google Drive
	ScopeSpreadsheetsReadonly  string = slides.SpreadsheetsReadonlyScope  // View your Google Spreadsheets
)
