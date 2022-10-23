package models

type Chat struct {
	EventUID    string `json:"event_uid,omitempty"`
	EventTitle  string `json:"event_title,omitempty"`
	EventAvatar string `json:"event_avatar,omitempty"`
}

type Message struct {
	UserID      int64  `json:"user_id,omitempty"`
	EventID     string `json:"event_id,omitempty"`
	Text        string `json:"text,omitempty"`
	TimeCreated int64  `json:"time_created,omitempty"`
}

type GetMessageOpts struct {
	EventID string
	Page    int
	Limit   int
}

//type MessageList []Message
//
//func (m MessageList) Len() int {
//	return len(m)
//}
//
//func (m MessageList) Less(i int, j int) bool {
//	return m[i].TimeCreated < m[j].TimeCreated
//}
//
//func (m MessageList) Swap(i int, j int) {
//	temp := m[j]
//	m[j] = m[i]
//	m[i] = temp
//}
