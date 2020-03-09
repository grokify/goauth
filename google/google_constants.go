package google

import "google.golang.org/api/slides/v1"

const (
	ClientSecretEnv         = "GOOGLE_APP_CLIENT_SECRET"
	EnvGoogleAppCredentials = "GOOGLE_APP_CREDENTIALS"
	EnvGoogleAppScopes      = "GOOGLE_APP_SCOPES"
	ScopeDrive              = slides.DriveScope                 // See, edit, create, and delete all of your Google Drive files
	ScopeDriveFile          = slides.DriveFileScope             // View and manage Google Drive files and folders that you have opened or created with this app
	ScopeDriveRead          = slides.DriveReadonlyScope         // See and download all your Google Drive files
	ScopePresentations      = slides.PresentationsScope         // View and manage your Google Slides presentations
	ScopePresentationsRead  = slides.PresentationsReadonlyScope // View your Google Slides presentations
	ScopeSpreadsheets       = slides.SpreadsheetsScope          // See, edit, create, and delete your spreadsheets in Google Drive
	ScopeSpreadsheetsRead   = slides.SpreadsheetsReadonlyScope  // View your Google Spreadsheets
	ScopeUserEmail          = "https://www.googleapis.com/auth/userinfo#email"
	ScopeUserProfile        = "https://www.googleapis.com/auth/userinfo.profile" // call https://www.googleapis.com/oauth2/v1/userinfo?alt=json

	// ScopeDrive             = "https://www.googleapis.com/auth/drive"
	// ScopeDriveRead         = "https://www.googleapis.com/auth/drive.readonly"
	// ScopeDriveFile         = "https://www.googleapis.com/auth/drive.file"
	// ScopeSheetsRead = "https://www.googleapis.com/auth/spreadsheets.readonly"
	// ScopeSheets     = "https://www.googleapis.com/auth/spreadsheets"
	// ScopeSlides = "https://www.googleapis.com/auth/presentations"
)

// DriveReadonlyDesc = "Allows read-only access to the user's file metadata and file content."
// DriveFileDesc     = "Per-file access to files created or opened by the app."
// DriveDesc         = "Full, permissive scope to access all of a user's files. Request this scope only when it is strictly necessary."
//SpreadsheetsReadonlyDesc = "Allows read-only access to the user's sheets and their properties."
//SpreadsheetsDesc         = "Allows read/write access to the user's sheets and their properties."
