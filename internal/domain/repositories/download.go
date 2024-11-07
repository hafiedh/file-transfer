package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"hafiedh.com/downloader/internal/domain/entities"
)

type (
	FileTransfer interface {
		SaveUpload(ctx context.Context, data *entities.MetaData) (err error)
		SaveDownload(ctx context.Context, data *entities.MetaData) (err error)
		FindByID(ctx context.Context, id int64) (data entities.MetaData, err error)
	}

	fileTransfer struct {
		db *pgxpool.Pool
	}
)

func NewFileTransfer(db *pgxpool.Pool) FileTransfer {
	return &fileTransfer{
		db: db,
	}
}

func (f *fileTransfer) SaveUpload(ctx context.Context, data *entities.MetaData) (err error) {
	_, err = f.db.Exec(ctx, `
		INSERT INTO meta_data_upload ( file_name, url, file_size, extension, status, created_at, updated_at, file_id )
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
	`, data.FileName, data.URL, data.FileSize, data.Extension, data.Status, data.CreatedAt, data.UpdatedAt, data.FileID)
	return
}

func (f *fileTransfer) SaveDownload(ctx context.Context, data *entities.MetaData) (err error) {
	_, err = f.db.Exec(ctx, `
		INSERT INTO meta_data_download ( file_name, url, file_size, extension, status, created_at, updated_at, file_id )
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
	`, data.FileName, data.URL, data.FileSize, data.Extension, data.Status, data.CreatedAt, data.UpdatedAt, data.FileID)
	return
}

func (f *fileTransfer) FindByID(ctx context.Context, id int64) (data entities.MetaData, err error) {
	err = f.db.QueryRow(ctx, `
		SELECT id, file_name, url, file_size, extension, status, created_at, updated_at
		FROM meta_data_upload
		WHERE id = $1
	`, id).Scan(&data.ID, &data.FileName, &data.URL, &data.FileSize, &data.Extension, &data.Status, &data.CreatedAt, &data.UpdatedAt)
	return
}
