package model

type Ticket struct {
	GuildId uint64 `json:"guild_id,string"`
	Id      int    `json:"id"`
}

type TicketStatusUpdate struct {
	Ticket
	ChannelId     uint64 `json:"channel_id,string"`
	NewCategoryId uint64 `json:"new_category_id,string"`
}
