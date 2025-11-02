CREATE TYPE listing_status AS ENUM (
    'draft',
    'pending_moderation',
    'published',
    'archived',
    'deleted',
    'rejected'
);

CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(15, 2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    category_id VARCHAR(100) NOT NULL,
    status listing_status NOT NULL DEFAULT 'draft',

    location JSONB NOT NULL DEFAULT '{}'::jsonb,
    
    metadata JSONB NOT NULL DEFAULT '{"responses_count": 0, "favorites_count": 0}'::jsonb,
    
    photo_ids TEXT[] DEFAULT '{}',
    video_ids TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',

    views_count BIGINT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT valid_currency CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT valid_price CHECK (price >= 0)
);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_listings_updated_at
    BEFORE UPDATE ON listings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
