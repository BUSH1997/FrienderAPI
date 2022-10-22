package models

type Chat struct {
	EventUID    string
	EventTitle  string
	EventAvatar string
}

type Message struct {
	UserID      int64
	EventID     string
	Text        string
	TimeCreated int64
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
