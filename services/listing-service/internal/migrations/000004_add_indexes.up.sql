-- Indexes for listings table
CREATE INDEX idx_listings_user_id ON listings(user_id);
CREATE INDEX idx_listings_status ON listings(status);
CREATE INDEX idx_listings_category_id ON listings(category_id);
CREATE INDEX idx_listings_created_at ON listings(created_at DESC);
CREATE INDEX idx_listings_published_at ON listings(published_at DESC) WHERE published_at IS NOT NULL;
CREATE INDEX idx_listings_price ON listings(price);

-- Composite indexes for common queries
CREATE INDEX idx_listings_status_created ON listings(status, created_at DESC);
CREATE INDEX idx_listings_user_status ON listings(user_id, status);
CREATE INDEX idx_listings_category_status ON listings(category_id, status);

-- GIN index for JSON fields
CREATE INDEX idx_listings_location_gin ON listings USING GIN (location);
CREATE INDEX idx_listings_metadata_gin ON listings USING GIN (metadata);

-- GIN index for arrays
CREATE INDEX idx_listings_tags_gin ON listings USING GIN (tags);

-- Full-text search index
CREATE INDEX idx_listings_title_fts ON listings USING GIN (to_tsvector('english', title));
CREATE INDEX idx_listings_description_fts ON listings USING GIN (to_tsvector('english', description));

-- Indexes for photos table
CREATE INDEX idx_photos_listing_id ON photos(listing_id);
CREATE INDEX idx_photos_order ON photos(listing_id, "order");

-- Indexes for videos table
CREATE INDEX idx_videos_listing_id ON videos(listing_id);

-- Partial index for active listings
CREATE INDEX idx_listings_active ON listings(created_at DESC) 
    WHERE status IN ('published', 'pending_moderation');
