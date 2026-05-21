package checkhttp

type CheckHTTP struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CheckGroup  string `json:"check_group,omitempty"`

	Enabled  int    `json:"enabled"`
	Place    string `json:"place"`
	Entities []int  `json:"entities"`

	Interval     int `json:"interval,omitempty"`
	CheckTimeout int `json:"check_timeout,omitempty"`

	TelegramChannelID int `json:"telegram_channel_id,omitempty"`
	SlackChannelID    int `json:"slack_channel_id,omitempty"`
	MMChannelID       int `json:"mm_channel_id,omitempty"`
	PDChannelID       int `json:"pd_channel_id,omitempty"`
	EmailChannelId    int `json:"email_channel_id,omitempty"`

	URL            string `json:"url"`
	Method         string `json:"method"`
	IgnoreSSLError int    `json:"ignore_ssl_error,omitempty"`

	AcceptedStatusCodes []int  `json:"accepted_status_codes"`
	Body                string `json:"body,omitempty"`
	BodyReq             string `json:"body_req,omitempty"`
	HeaderReq           string `json:"header_req,omitempty"`

	Retries   int    `json:"retries,omitempty"`
	Redirects int    `json:"redirects,omitempty"`
	Runbook   string `json:"runbook,omitempty"`
}

type (
	Request  = CheckHTTP
	Response = CheckHTTP
)
