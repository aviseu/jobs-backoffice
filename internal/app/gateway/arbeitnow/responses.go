package arbeitnow

import (
	"gopkg.in/guregu/null.v3"
)

type jobBoardResponse struct {
	Jobs  []*jobEntry `json:"data"`
	Links struct {
		Next null.String `json:"next"`
	} `json:"links"`
}

type jobEntry struct {
	Slug        string   `json:"slug"`
	CompanyName string   `json:"company_name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Location    string   `json:"location"`
	Tags        []string `json:"tags"`
	JobTypes    []string `json:"job_types"`
	CreatedAt   int64    `json:"created_at"`
	Remote      bool     `json:"remote"`
}
