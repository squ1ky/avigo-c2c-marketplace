-- Drop indexes for listings
DROP INDEX IF EXISTS idx_listings_user_id;
DROP INDEX IF EXISTS idx_listings_status;
DROP INDEX IF EXISTS idx_listings_category_id;
DROP INDEX IF EXISTS idx_listings_created_at;
DROP INDEX IF EXISTS idx_listings_published_at;
DROP INDEX IF EXISTS idx_listings_price;
DROP INDEX IF EXISTS idx_listings_status_created;
DROP INDEX IF EXISTS idx_listings_user_status;
DROP INDEX IF EXISTS idx_listings_category_status;
DROP INDEX IF EXISTS idx_listings_location_gin;
DROP INDEX IF EXISTS idx_listings_metadata_gin;
DROP INDEX IF EXISTS idx_listings_tags_gin;
DROP INDEX IF EXISTS idx_listings_title_fts;
DROP INDEX IF EXISTS idx_listings_description_fts;
DROP INDEX IF EXISTS idx_listings_active;

-- Drop indexes for photos
DROP INDEX IF EXISTS idx_photos_listing_id;
DROP INDEX IF EXISTS idx_photos_order;

-- Drop indexes for videos
DROP INDEX IF EXISTS idx_videos_listing_id;
