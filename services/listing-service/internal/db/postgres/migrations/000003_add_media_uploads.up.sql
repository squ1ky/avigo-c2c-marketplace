CREATE TABLE media_uploads
(
    id         UUID PRIMARY KEY,
    user_id    UUID         NOT NULL,
    s3_key     VARCHAR(512) NOT NULL,
    filename   VARCHAR(255),
    mime_type  VARCHAR(100),
    size       BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_media_uploads_user_id ON media_uploads (user_id);