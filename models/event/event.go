package event

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/validation"
)

type Event struct {
	Id          int64   `json:"id" db:"id"`
	OwnerId     int64   `json:"owner_id" db:"owner_id"`
	Description string  `json:"description" db:"description"` // omit
	Title       string  `json:"title" db:"title"`
	Active      int     `json:"active" db:"is_active"`
	Date        string  `json:"date" db:"date"`
	Cost        float64 `json:"cost" db:"cost"`
	Attendees   []int64 `json:"attendees" db:"-"`
	CreatedAt   string  `json:"created_at,omitempty" db:"-"`
	UpdatedAt   string  `json:"updated_at,omitempty" db:"-"`
}

func (e *Event) Save() (int64, error) {
	columns, values := e.getColumnStringAndValues()

	id, err := EventInsert(columns, values)
	e.Id = id

	return id, err
}

/**
 * Build a new user from a user map
 */
func (e *Event) Build(eventPostData map[string][]interface{}) error {
	var err error

	for k, v := range eventPostData {
		switch k {
		case "title":
			title := v[0].(string)

			if !validation.IsNonEmptyString(title) {
				err = errors.New(messaging.EVENT_TITLE_EMPTY)
			} else {
				e.Title = title
			}
		case "date":
			// TODO: validate date
			date := v[0].(string)
			e.Date = date
		case "owner_id":
			owner_id := v[0].(int64)

			if owner_id <= 0 {
				err = errors.New(messaging.EVENT_OWNER_ID_INVALID)
			} else {
				e.OwnerId = owner_id
			}
		case "description":
			e.Description = v[0].(string)
		case "cost":
			e.Cost = v[0].(float64)
		default:
			// only add properties that are user supplied,
			// all else can be ignored here
		}

		// return err as soon as one arises
		if err != nil {
			break
		}
	}

	return err
}

func (e *Event) getColumnStringAndValues() (string, []interface{}) {
	var columns string
	var values []interface{}

	eventElem := reflect.ValueOf(e).Elem()
	eventType := eventElem.Type()

	// Iterate over fields
	for i := 0; i < eventElem.NumField(); i++ {
		columnName := eventType.Field(i).Tag.Get("db")

		if columnName == "id" || columnName == "-" {
			continue
		}

		columns += fmt.Sprintf("%s, ", columnName)
		values = append(values, eventElem.Field(i).Interface())
	}

	columns = strings.TrimRight(columns, ", ")

	return columns, values
}
