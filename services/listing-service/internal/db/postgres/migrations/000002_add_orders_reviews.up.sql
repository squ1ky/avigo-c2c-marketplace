ALTER TABLE listings
    ADD COLUMN is_sold BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE orders
(
    id            UUID PRIMARY KEY,
    listing_id    UUID        NOT NULL REFERENCES listings (id),
    buyer_id      UUID        NOT NULL,
    seller_id     UUID        NOT NULL,
    status        VARCHAR(50) NOT NULL     DEFAULT 'created', -- 'created', 'completed', 'cancelled'
    cancel_reason TEXT,

    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_orders_buyer_id ON orders (buyer_id);
CREATE INDEX idx_orders_seller_id ON orders (seller_id);
CREATE INDEX idx_orders_listing_id ON orders (listing_id);


CREATE TABLE reviews
(
    id               UUID PRIMARY KEY,
    order_id         UUID    NOT NULL REFERENCES orders (id),
    reviewer_id      UUID    NOT NULL, -- Buyer
    reviewed_user_id UUID    NOT NULL, -- Seller
    rating           INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    text             TEXT,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reviews_reviewed_user_id ON reviews (reviewed_user_id);
CREATE UNIQUE INDEX idx_reviews_order_id ON reviews (order_id);
