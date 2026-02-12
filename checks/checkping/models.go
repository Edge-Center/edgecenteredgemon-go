package checkping

type CheckPing struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CheckGroup  string `json:"check_group,omitempty"`

	Enabled  int    `json:"enabled"`
	Place    string `json:"place"`
	Entities []int  `json:"entities"`

	PacketSize   int `json:"packet_size,omitempty"`
	CountPackets int `json:"count_packets,omitempty"`

	Interval     int `json:"interval,omitempty"`
	CheckTimeout int `json:"check_timeout"`

	TelegramChannelID int `json:"telegram_channel_id,omitempty"`
	SlackChannelID    int `json:"slack_channel_id,omitempty"`
	MMChannelID       int `json:"mm_channel_id,omitempty"`
	PDChannelID       int `json:"pd_channel_id,omitempty"`
	EmailChannelId    int `json:"email_channel_id,omitempty"`

	IP      string `json:"ip"`
	Retries int    `json:"retries,omitempty"`
	Runbook string `json:"runbook,omitempty"`
}

type Request = CheckPing
type Response = CheckPing

type CreateResponse struct {
	ID int `json:"id"`
}
