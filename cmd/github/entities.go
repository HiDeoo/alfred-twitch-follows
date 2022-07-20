package main

type GHRepo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	HtmlURL  string `json:"html_url"`
	PushedAt string `json:"pushed_at"`
}

type GHError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}
