package xmatters

type Shift struct {
	ID         *string          `json:"id"`
	Group      *GroupReference  `json:"group"`
	Name       *string          `json:"name"`
	Start      *string          `json:"start"`
	End        *string          `json:"end"`
	Timezone   *string          `json:"timezone"`
	Recurrence *ShiftRecurrence `json:"recurrence"`
	Members    []*ShiftMember   `json:"members"`
}

type ShiftPagination struct {
	*Pagination
	Shifts []*Shift `json:"data"`
}

type ShiftRecurrence struct {
	Frequency           *string   `json:"frequency"`
	RepeatEvery         *int64    `json:"repeatEvery,omitempty"`
	OnDays              []*string `json:"onDays,omitempty"`
	On                  *string   `json:"on,omitempty"`
	Months              []*string `json:"months,omitempty"`
	DateOfMonth         *string   `json:"dateOfMonth,omitempty"`
	DayOfWeekClassifier *string   `json:"dayOfWeekClassifier,omitempty"`
	DayOfWeek           *string   `json:"dayOfWeek,omitempty"`
	End                 *ShiftEnd `json:"end,omitempty"`
}

type ShiftEnd struct {
	EndBy       *string `json:"endBy"`
	Date        *string `json:"date"`
	Repetitions *int64  `json:"repititions"`
}

type ShiftMember struct {
	Recipient      *RecipientPointer `json:"recipient"`
	Shift          *ReferenceById    `json:"shift"`
	Position       *int64            `json:"position"`
	Delay          *int64            `json:"delay"`
	EscalationType *string           `json:"escalationType"`
	InRotation     *bool             `json:"inRotation"`
}

// RecipientPointer is a reference to a recipient.
type RecipientPointer struct {
	ID   *string `json:"id"`
	Type *string `json:"recipientType"`
}
