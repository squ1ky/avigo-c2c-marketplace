-- Drop trigger
DROP TRIGGER IF EXISTS update_listings_updated_at ON listings;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop table
DROP TABLE IF EXISTS listings;

-- Drop enum type
DROP TYPE IF EXISTS listing_status;
