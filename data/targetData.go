package data

type TargetData []struct {
	Category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Organization struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"organization"`
	Codename  string `json:"codename"`
	Slug      string `json:"slug"`
	OutageWin []struct {
		StartDate   int `json:"start_date"`
		EndDate     int `json:"end_date"`
		OutageStart int `json:"outage_starts_on"`
		OutageEnds  int `json:"outage_ends_on"`
		Options     struct {
			Days      []int  `json:"days"`
			Frequency string `json:"frequency"`
		} `json:"options"`
		WindowActive bool `json:"is_window_active"`
	} `json:"outage_windows"`
	SRT_Notes   string   `json:"srt_notes"`
	DateUpdated int      `json:"dateUpdated"`
	Active      bool     `json:"isActive"`
	New         bool     `json:"isNew"`
	Registered  bool     `json:"isRegistered"`
	Name        string   `json:"name"`
	AvgPayout   float64  `json:"averagePayout"`
	LastSubm    int      `json:"lastSubmitted"`
	VulnDisc    bool     `json:"vulnerability_discovery"`
	Workspace   bool     `json:"workspace_access_missing"`
	Updated     bool     `json:"isUpdated"`
	Incentives  []string `json:"incentives"`
}
