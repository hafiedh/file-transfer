package entities

import "time"

type MetaData struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	FileID    string    `db:"file_id"`
	FileName  string    `db:"file_name"`
	URL       string    `db:"url"`
	Extension string    `db:"extension"`
	Status    string    `db:"status"`
	ID        int64     `db:"id"`
	FileSize  int64     `db:"file_size"`
}
