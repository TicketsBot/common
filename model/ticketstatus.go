package model

type TicketStatus string

const (
	TicketStatusOpen    TicketStatus = "OPEN"
	TicketStatusPending TicketStatus = "PENDING"
	TicketStatusClosed  TicketStatus = "CLOSED"
)
