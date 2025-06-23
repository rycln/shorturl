package models

// Stats represents statistics for a URL shortener service.
// It contains aggregate counts of URLs and users in the system.
type Stats struct {
	// URLs is the total number of shortened URLs in the service
	URLs int `json:"urls"`

	// Users is the total number of registered users in the service
	Users int `json:"users"`
}
