package data

type MissionData []struct {
	TaskTemplate struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
	} `json:"taskTemplate"`
	PublishedOn  string   `json:"publishedOn"`
	ValidatedOn  string   `json:"validatedOn"`
	AttackType   []string `json:"attackType"`
	Organization struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"organization"`
	DeletedOn          string `json:"deletedOn"`
	StructuredResponse string `json:"structuredResponse"`
	ResponseType       string `json:"responseType"`
	ValidResponses     []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"validResponses"`
	CompletedOn     string   `json:"completedOn"`
	CreatedOn       string   `json:"createdOn"`
	Response        string   `json:"response"`
	Assignee        string   `json:"asignee"`
	TaskGroup       string   `json:"task_group"`
	PausedDuration  int      `json:"pausedDurationInSecs"`
	ReturnedForEdit int      `json:"returnedForEditOn"`
	Reviewer        string   `json:"reviewer"`
	ClaimedOn       string   `json:"claimedOn"`
	DeactivatedOn   string   `json:"deactivatedOn"`
	InvalidationOn  string   `json:"invalidatedOn"`
	TaskType        string   `json:"taskType"`
	AssetType       []string `json:"assetType"`
	AssigneeUser    string   `json:"asigneeUser"`
	Payout          struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	} `json:"payout"`
	Position              int    `json:"position"`
	Scope                 string `json:"scope"`
	CreatedBy             string `json:"createdBy"`
	ReviewedOn            string `json:"reviewedOn"`
	IsAssigneeCurrentUser bool   `json:"isAssigneeCurrentUser"`
	ID                    string `json:"id"`
	PausedOn              string `json:"pausedOn"`
	Blank                 int    `json:""`
	Title                 string `json:"title"`
	CampaignTempl         struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
	} `json:"campaignTemplate"`
	CanEditResponse bool     `json:"canEditResponse"`
	SV              []string `json:"sv"`
	Relations       struct{} `json:"relations"`
	Claims          struct {
		Current int `json:"current"`
		Limit   int `json:"limit"`
	} `json:"claims"`
	Attempts      int      `json:"attempts"`
	ModifiedOn    string   `json:"modifiedOn"`
	ReviewerUser  string   `json:"reviewerUser"`
	HasBeenViewed bool     `json:"hasBeenViewed"`
	CWE           []string `json:"cwe"`
	Description   string   `json:"description"`
	Campaign      struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"campaign"`
	BatchID         string `json:"batch_id"`
	SubmissionCount int    `json:"submissionsCount"`
	Listing         struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"listing"`
	Definition struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"definition"`
	Categories   []string `json:"categories"`
	DurationSecs int      `json:"durationInSecs"`
	ModifiedBy   string   `json:"modifiedBy"`
	Version      int      `json:"version"`
	Unauthorized string   `json:"unauthorizedAssignees"`
}
