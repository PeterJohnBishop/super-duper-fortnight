package clkup

// response types

type UserResponse struct {
	User User `json:"user"`
}

type TeamsResponse struct {
	Teams []Workspace `json:"teams"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}

type TasksResponse struct {
	Task []Task `json:"tasks"`
}

// hierarchy types

type Workspace struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Avatar  string   `json:"avatar"`
	Members []Member `json:"members"`
}

type Space struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Private           bool          `json:"private"`
	Color             string        `json:"color"`
	Avatar            string        `json:"avatar"`
	AdminCanManage    bool          `json:"admin_can_manage"`
	Archived          bool          `json:"archived"`
	Members           []SpaceMember `json:"members"`
	Statuses          []Status      `json:"statuses"`
	MultipleAssignees bool          `json:"multiple_assignees"`
	Features          Features      `json:"features"`
}

type Folder struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Orderindex       int64         `json:"orderindex"`
	OverrideStatuses bool          `json:"override_statuses"`
	Hidden           bool          `json:"hidden"`
	Space            SpaceLocation `json:"space"`
	TaskCount        string        `json:"task_count"`
	Lists            []interface{} `json:"lists"`
}

type List struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Orderindex       int64          `json:"orderindex"`
	Content          string         `json:"content"`
	Status           Status         `json:"status"`
	Priority         Priority       `json:"priority"`
	Assignee         string         `json:"assignee"`
	TaskCount        string         `json:"task_count"`
	DueDate          string         `json:"due_date"`
	StartDate        string         `json:"start_date"`
	Folder           FolderLocation `json:"folder"`
	Space            SpaceLocation  `json:"space"`
	Archived         bool           `json:"archived"`
	OverrideStatuses bool           `json:"override_statuses"`
	PermissionLevel  string         `json:"permission_level"`
}

type Task struct {
	Points          any            `json:"points"`
	DateDone        any            `json:"date_done"`
	GroupAssignees  []any          `json:"group_assignees"`
	Id              string         `json:"id"`
	Orderindex      string         `json:"orderindex"`
	Parent          any            `json:"parent"`
	Sharing         Sharing        `json:"sharing"`
	TextContent     string         `json:"text_content"`
	TopLevelParent  any            `json:"top_level_parent"`
	Assignees       []any          `json:"assignees"`
	Space           SpaceLocation  `json:"space"`
	StartDate       any            `json:"start_date"`
	Tags            []any          `json:"tags"`
	TeamId          string         `json:"team_id"`
	Url             string         `json:"url"`
	Watchers        []Watcher      `json:"watchers"`
	CustomFields    []CustomField  `json:"custom_fields"`
	CustomItemId    int            `json:"custom_item_id"`
	LinkedTasks     []any          `json:"linked_tasks"`
	Project         Project        `json:"project"`
	Status          Status         `json:"status"`
	Checklists      []any          `json:"checklists"`
	DateUpdated     string         `json:"date_updated"`
	Dependencies    []any          `json:"dependencies"`
	Description     string         `json:"description"`
	List            ListLocation   `json:"list"`
	PermissionLevel string         `json:"permission_level"`
	Priority        any            `json:"priority"`
	TimeSpent       int            `json:"time_spent"`
	Name            string         `json:"name"`
	TimeEstimate    any            `json:"time_estimate"`
	Attachments     []any          `json:"attachments"`
	Creator         Creator        `json:"creator"`
	DateCreated     string         `json:"date_created"`
	Locations       []any          `json:"locations"`
	Archived        bool           `json:"archived"`
	DateClosed      any            `json:"date_closed"`
	CustomId        any            `json:"custom_id"`
	DueDate         any            `json:"due_date"`
	Folder          FolderLocation `json:"folder"`
}

// user types

type User struct {
	ID                int64  `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Color             string `json:"color"`
	ProfilePicture    string `json:"profilePicture"`
	Initials          string `json:"initials"`
	WeekStartDay      int64  `json:"week_start_day"`
	GlobalFontSupport bool   `json:"global_font_support"`
	Timezone          string `json:"timezone"`
}

type Member struct {
	User WorkspaceUser `json:"user"`
}

type WorkspaceUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
}

type SpaceMember struct {
	User SpaceUser `json:"user"`
}

type SpaceUser struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
	Initials       string `json:"initials"`
}

type Watcher struct {
	Id             int    `json:"id"`
	Initials       string `json:"initials"`
	ProfilePicture any    `json:"profilePicture"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	Email          string `json:"email"`
}

// property types

type SpaceStatus struct {
	Status     string `json:"status"`
	Type       string `json:"type"`
	OrderIndex int    `json:"orderindex"`
	Color      string `json:"color"`
}

type Features struct {
	DueDates          DueDatesFeature `json:"due_dates"`
	TimeTracking      FeatureEnabled  `json:"time_tracking"`
	Tags              FeatureEnabled  `json:"tags"`
	TimeEstimates     FeatureEnabled  `json:"time_estimates"`
	Checklists        FeatureEnabled  `json:"checklists"`
	CustomFields      FeatureEnabled  `json:"custom_fields"`
	RemapDependencies FeatureEnabled  `json:"remap_dependencies"`
	DependencyWarning FeatureEnabled  `json:"dependency_warning"`
	Portfolios        FeatureEnabled  `json:"portfolios"`
}

type FeatureEnabled struct {
	Enabled bool `json:"enabled"`
}

type DueDatesFeature struct {
	Enabled            bool `json:"enabled"`
	StartDate          bool `json:"start_date"`
	RemapDueDates      bool `json:"remap_due_dates"`
	RemapClosedDueDate bool `json:"remap_closed_due_date"`
}

type SpaceLocation struct {
	Id string `json:"id"`
}

type Status struct {
	Orderindex int    `json:"orderindex"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Color      string `json:"color"`
	Id         string `json:"id"`
}

type ListLocation struct {
	Access bool   `json:"access"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

type Priority struct {
	Priority string `json:"priority"`
	Color    string `json:"color"`
}

type Sharing struct {
	Public               bool     `json:"public"`
	PublicFields         []string `json:"public_fields"`
	PublicShareExpiresOn any      `json:"public_share_expires_on"`
	SeoOptimized         bool     `json:"seo_optimized"`
	Token                any      `json:"token"`
}

type CustomField struct {
	Name           string     `json:"name"`
	Required       bool       `json:"required"`
	Type           string     `json:"type"`
	TypeConfig     TypeConfig `json:"type_config"`
	DateCreated    string     `json:"date_created"`
	HideFromGuests bool       `json:"hide_from_guests"`
	Id             string     `json:"id"`
}

type Project struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
}

type Creator struct {
	Color          string `json:"color"`
	Email          string `json:"email"`
	Id             int    `json:"id"`
	ProfilePicture any    `json:"profilePicture"`
	Username       string `json:"username"`
}

type FolderLocation struct {
	Name   string `json:"name"`
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
	Id     string `json:"id"`
}

type TypeConfig struct{}

type Performance struct {
	Duration string
	RPM      string
	TPS      string
}
