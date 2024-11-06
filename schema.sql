CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS meta_data_upload (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    file_name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS meta_data_upload_file_name_index ON meta_data_upload USING btree (file_name);
CREATE INDEX IF NOT EXISTS meta_data_upload_url_index ON meta_data_upload USING btree (url);
CREATE INDEX IF NOT EXISTS meta_data_upload_user_id_index ON meta_data_upload USING btree (user_id);

CREATE TABLE IF NOT EXISTS meta_data_download (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    file_name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    extension VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    file_size INT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS meta_data_download_file_name_index ON meta_data_download USING btree (file_name);
CREATE INDEX IF NOT EXISTS meta_data_download_url_index ON meta_data_download USING btree (url);
CREATE INDEX IF NOT EXISTS meta_data_download_user_id_index ON meta_data_download USING btree (user_id);


