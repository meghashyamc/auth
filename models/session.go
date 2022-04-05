package models

type Session struct {
	LoggedIn               bool  `json:"logged_in"`
	LatestInvalidIssueTime int64 `json:"latest_invalid_issue_time"`
}
