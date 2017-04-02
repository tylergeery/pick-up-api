package event

import (
	"errors"
	"fmt"
	"log"

	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/resources"
)

/**
 * Get an Event Object from the DB
 */
func EventGetById(id int64) (Event, error) {
	var event Event
	tx := resources.TX()

	rows, err := tx.Query("SELECT id, owner_id, title, description, date, cost, is_active, created_at, updated_at FROM events WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id, owner_id int64
		var title, description, date, created_at, updated_at string
		var cost float64
		var is_active int

		if err := rows.Scan(&id, &owner_id, &title, &description, &date, &cost, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}

		event = Event{
			Id:          id,
			OwnerId:     owner_id,
			Description: description,
			Title:       title,
			Date:        date,
			Cost:        cost,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
			Active:      is_active,
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if event.Id == 0 && err == nil {
		err = errors.New(messaging.EVENT_NOT_FOUND)
	}

	if err == nil && event.Active != 1 {
		err = errors.New(messaging.EVENT_NOT_ACTIVE)
	}

	return event, err
}

/**
 * Insert an event to the database
 */
func EventInsert(columns string, values []interface{}) (int64, error) {
	var id int64

	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   INSERT INTO events (%s, created_at)
            VALUES (%s, NOW()) RETURNING id
        `, columns, resources.SqlStub(len(values)))

	err := tx.QueryRow(query, values...).Scan(&id)

	if err == nil {
		tx.Commit()
	}

	return id, err
}

/**
 * Update an event in the DB
 */
func EventUpdateValues(eventId int64, columns string, values []interface{}) error {
	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   UPDATE events
            SET (%s, updated_at) = (%s, NOW())
            WHERE id = %d
        `, columns, resources.SqlStub(len(values)), eventId)

	_, err := tx.Exec(query, values...)

	if err == nil {
		tx.Commit()
	}

	return err
}

/**
 * Mark an event as inactive in the DB
 */
func EventDelete(eventId int64) error {
	tx := resources.TX()
	defer tx.Rollback()
	stmt, err := tx.Prepare(
		`   UPDATE events
            SET is_active = 0, updated_at = NOW()
            WHERE id = $1;`)

	if err == nil {
		defer stmt.Close()
		_, err = stmt.Exec(eventId)
	}

	if err == nil {
		err = tx.Commit()
	}

	return err
}
