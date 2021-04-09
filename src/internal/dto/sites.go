package dto

import (
	"fmt"
	"strings"
)

type Site struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	BaseURL    string `db:"base_url" json:"base_url"`
	SearchURL  string `db:"search_url" json:"search_url"`
	BasePath   string `db:"base_path" json:"base_path"`
	UseSubPath bool   `db:"use_sub_path" json:"use_sub_path"`
	MetaData   string `db:"meta_data" json:"meta_data"`
	MinAge     int    `db:"min_age" json:"min_age"`
}

func (s *Site) GetURL(search map[string]string) string {
	output := s.SearchURL
	for k, v := range search {
		output = strings.ReplaceAll(output, fmt.Sprintf("{%s}", k), v)
	}
	return output
}
