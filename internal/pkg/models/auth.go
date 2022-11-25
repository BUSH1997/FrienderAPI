package models

type AuthToken struct {
	UserID  int64
	Value   string
	Expires int64
}

type RefreshToken struct {
	Value       string
	Expires     int64
	FingerPrint string
}

type FingerPrintData struct {
	UserAgent string
	UserIP    string
}

type AuthParams struct {
	UserID int64  `json:"user_id,omitempty"`
	AppID  int64  `json:"app_id,omitempty"`
	Time   int64  `json:"time,omitempty"`
	Sign   string `json:"sign,omitempty"`
}
