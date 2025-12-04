DROP INDEX IF EXISTS idx_reviews_order_id;
DROP INDEX IF EXISTS idx_reviews_reviewed_user_id;
DROP TABLE IF EXISTS reviews;

DROP INDEX IF EXISTS idx_orders_listing_id;
DROP INDEX IF EXISTS idx_orders_seller_id;
DROP INDEX IF EXISTS idx_orders_buyer_id;
DROP TABLE IF EXISTS orders;

ALTER TABLE listings DROP COLUMN IF EXISTS is_sold;