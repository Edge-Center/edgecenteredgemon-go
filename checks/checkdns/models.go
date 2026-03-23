package checkdns

type CheckDNS struct {
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

	IP       string `json:"ip"`
	Port     int    `json:"port,omitempty"`
	Resolver string `json:"resolver"`

	RecordType string `json:"record_type"`
	Retries    int    `json:"retries,omitempty"`
	Runbook    string `json:"runbook,omitempty"`
}

type (
	Request  = CheckDNS
	Response = CheckDNS
)

type CreateResponse struct {
	ID int `json:"id"`
}
