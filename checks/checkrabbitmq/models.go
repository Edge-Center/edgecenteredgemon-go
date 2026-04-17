package checkrabbitmq

type CheckRabbitMQ struct {
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
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`

	Vhost   string `json:"vhost"`
	Retries int    `json:"retries,omitempty"`
	Runbook string `json:"runbook,omitempty"`
}

type (
	Request  = CheckRabbitMQ
	Response = CheckRabbitMQ
)
