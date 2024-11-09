package model

type TicketStatusUpdate struct {
	Ticket
	ChannelId     uint64 `json:"channel_id,string"`
	NewCategoryId uint64 `json:"new_category_id,string"`
}
