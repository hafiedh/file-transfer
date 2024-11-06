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
		INSERT INTO meta_data_upload ( file_name, url, file_size, extension, status, created_at, updated_at )
		VALUES ( $1, $2, $3, $4, $5, $6, $7 )
	`, data.FileName, data.URL, data.FileSize, data.Extension, data.Status, data.CreatedAt, data.UpdatedAt)
	return
}

func (f *fileTransfer) SaveDownload(ctx context.Context, data *entities.MetaData) (err error) {
	_, err = f.db.Exec(ctx, `
		INSERT INTO meta_data_download ( file_name, url, file_size, extension, status, created_at, updated_at )
		VALUES ( $1, $2, $3, $4, $5, $6, $7 )
	`, data.FileName, data.URL, data.FileSize, data.Extension, data.Status, data.CreatedAt, data.UpdatedAt)
	return
}
