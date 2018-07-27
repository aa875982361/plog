package model

type WebLog struct {
	Type       string `json:"type"`
	Project    string `json:"project"`
	User       string `json:"user"`
	Tag        string `json:"tag"`
	Detail     string `json:"detail"`
	CreateTime int64  `json:"createtime"`
}

func (w *WebLog) Insert(openID string) error {

	return nil
}
