package entities

import "time"

type MetaData struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	FileName  string    `db:"file_name"`
	URL       string    `db:"url"`
	Extension string    `db:"extension"`
	Status    string    `db:"status"`
	ID        int       `db:"id"`
	FileSize  int       `db:"file_size"`
}
