package models

import (
	"time"
)

const (
	StatusDraft     string = "draft"
	StatusPublished string = "published"
	StatusOngoing   string = "ongoing"
	StatusCompleted string = "completed"
	StatusCancelled string = "cancelled"
	StatusPostponed string = "postponed"
)

type EventUpdateRequest struct {
	Id           int    `validate:"required"`
	Title        string `validate:"min=5,max=255"`
	About        string `validate:"min=5"`
	StartDate    time.Time
	Location     string
	Status       string `validate:"oneof=draft published ongoing completed cancelled postponed"`
	MaxAttendees int    `validate:"min=5,max=1000"`
}

type EventCreateRequest struct {
	Title        string    `validate:"required,min=5,max=255"`
	About        string    `validate:"required,min=5"`
	StartDate    time.Time `validate:"required"`
	Location     string    `validate:"required"`
	Status       string    `validate:"required,oneof=draft published ongoing completed cancelled postponed"`
	MaxAttendees int       `validate:"required,min=5,max=1000"`
	Creator      string    `validate:"required,uuid"`
}

type EventResponse struct {
	Id                int
	Title             string
	About             string
	StartDate         time.Time
	Location          string
	Status            string
	MaxAttendees      int
	CurrentAttendance int
	Creator           string
}