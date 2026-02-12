package statuspage

type Base struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	CustomStyle string `json:"custom_style,omitempty"`
}
type Request struct {
	Base
	Checks []int `json:"checks"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

type Response struct {
	ID int `json:"id"`
	Base
	Checks []Checks `json:"checks"`
}

type Checks struct {
	CheckID int `json:"multi_check_id"`
}
