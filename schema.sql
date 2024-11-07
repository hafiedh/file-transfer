CREATE TABLE IF NOT EXISTS meta_data_upload (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    file_id VARCHAR(255) DEFAULT NULL,
    file_name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    file_size INT NOT NULL
);

CREATE INDEX IF NOT EXISTS meta_data_upload_file_name_index ON meta_data_upload USING btree (file_name);
CREATE INDEX IF NOT EXISTS meta_data_upload_url_index ON meta_data_upload USING btree (url);

CREATE TABLE IF NOT EXISTS meta_data_download (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    file_id VARCHAR(255) DEFAULT NULL,
    file_name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    file_size INT NOT NULL
);

CREATE INDEX IF NOT EXISTS meta_data_download_file_name_index ON meta_data_download USING btree (file_name);
CREATE INDEX IF NOT EXISTS meta_data_download_url_index ON meta_data_download USING btree (url);


