// mirror tables in schema.sql using structs. 
// Use standard types and add metadata tags ( db: or json: ) to tell db libraries and JSON encoders how to map the fields.

package repository

import (
	"time"
)

// BoardModel maps directly to the boards table in schema.sql
type BoardModel struct {
	ID        string    `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// ColumnModel maps directly to the columns table in schema.sql
type ColumnModel struct {
	ID		string    `db:"id" json:"id"`
	BoardID	string    `db:"board_id" json:"boardId"`
	Title	string    `db:"title" json:"title"`
	Position	int       `db:"position" json:"position"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// TaskModel maps directly to the tasks table in schema.sql
type TaskModel struct {
	ID          string    `db:"id" json:"id"`
	ColumnID    string    `db:"column_id" json:"columnId"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    int       `db:"position" json:"position"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Representing clean tree structure in domain layer that matches the exact tree structure the UI needs.
type BoardAggregate struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Columns []ColumnAggregate `json:"columns"`
}

type ColumnAggregate struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Position int `json:"position"`
	Tasks []TaskModel `json:"tasks"`
}