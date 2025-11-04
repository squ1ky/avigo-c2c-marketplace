CREATE TABLE IF NOT EXISTS videos (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL,
    object_id VARCHAR(255) NOT NULL, -- MinIO/S3 object ID
    url TEXT NOT NULL,               -- Public URL or presigned URL
    thumbnail_id VARCHAR(255),       -- ID of thumbnail image
    duration INTEGER NOT NULL DEFAULT 0 CHECK (duration >= 0),
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL CHECK (size > 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

                             CONSTRAINT fk_videos_listing
                             FOREIGN KEY (listing_id)
    REFERENCES listings(id)
                         ON DELETE CASCADE,

    CONSTRAINT valid_duration CHECK (duration >= 0),
    CONSTRAINT valid_size CHECK (size > 0),
    CONSTRAINT valid_content_type CHECK (content_type LIKE 'video/%')
    );
