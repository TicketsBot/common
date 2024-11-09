package model

type Ticket struct {
	GuildId uint64 `json:"guild_id,string"`
	Id      int    `json:"id"`
}
