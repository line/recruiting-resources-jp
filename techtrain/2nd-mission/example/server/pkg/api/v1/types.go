package v1

import (
	"time"
)

type JsonTime struct {
	time.Time
}

// use YY-mm-DDTHH:MM:SSZ time format
func (j JsonTime) format() string {
	return j.Time.Format("2006-01-02T15:04:05Z")
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

type TodoResponse struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Finished    bool     `json:"finished"`
	CreatedAt   JsonTime `json:"createdAt"`
	UpdatedAt   JsonTime `json:"updatedAt"`
}

type TodoListResponse struct {
	List []TodoResponse `json:"list"`
}

type TodoCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Finished    *bool   `json:"finished"`
}
