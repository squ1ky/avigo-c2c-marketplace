-- ===================================================================================
-- 1. ЭЛЕКТРОНИКА (Level 0)
-- ID: d290f1ee-6c54-4b01-90e6-d701748f0851
-- ===================================================================================
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('d290f1ee-6c54-4b01-90e6-d701748f0851', 'Электроника', 'electronics', NULL, 0, 100, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 1.1 Телефоны (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('e1ec1111-1111-1111-1111-111111111111', 'Телефоны', 'phones', 'd290f1ee-6c54-4b01-90e6-d701748f0851', 1, 110, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 1.1.1 Смартфоны (Level 2)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('e1ec2222-2222-2222-2222-222222222222', 'Смартфоны', 'smartphones', 'e1ec1111-1111-1111-1111-111111111111', 2, 111, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 1.1.2 Аксессуары для телефонов (Level 2)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('e1ec3333-3333-3333-3333-333333333333', 'Аксессуары', 'phone-accessories', 'e1ec1111-1111-1111-1111-111111111111', 2, 112, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 1.2 Ноутбуки (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('e1ec4444-4444-4444-4444-444444444444', 'Ноутбуки', 'laptops', 'd290f1ee-6c54-4b01-90e6-d701748f0851', 1, 120, NOW())
    ON CONFLICT (id) DO NOTHING;


-- ===================================================================================
-- 2. ТРАНСПОРТ (Level 0)
-- ID: a0ee5555-5555-5555-5555-555555555555
-- ===================================================================================
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('a0ee5555-5555-5555-5555-555555555555', 'Транспорт', 'transport', NULL, 0, 200, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 2.1 Автомобили (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('a0ee6666-6666-6666-6666-666666666666', 'Автомобили', 'cars', 'a0ee5555-5555-5555-5555-555555555555', 1, 210, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 2.2 Запчасти (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('a0ee7777-7777-7777-7777-777777777777', 'Запчасти', 'auto-parts', 'a0ee5555-5555-5555-5555-555555555555', 1, 220, NOW())
    ON CONFLICT (id) DO NOTHING;


-- ===================================================================================
-- 3. ОДЕЖДА И ОБУВЬ (Level 0)
-- ID: c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
-- ===================================================================================
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Одежда и обувь', 'clothes-shoes', NULL, 0, 300, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 3.1 Женская одежда (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('c300bbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Женская одежда', 'womens-clothing', 'c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, 310, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 3.1.1 Платья (Level 2)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('c300cccc-cccc-cccc-cccc-cccccccccccc', 'Платья и юбки', 'dresses', 'c300bbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 2, 311, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 3.2 Мужская обувь (Level 1)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('c300dddd-dddd-dddd-dddd-dddddddddddd', 'Мужская обувь', 'mens-shoes', 'c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, 320, NOW())
    ON CONFLICT (id) DO NOTHING;

-- 3.2.1 Кроссовки (Level 2)
INSERT INTO categories (id, name, slug, parent_id, level, "order", created_at)
VALUES
    ('c300eeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Кроссовки', 'sneakers', 'c300dddd-dddd-dddd-dddd-dddddddddddd', 2, 321, NOW())
    ON CONFLICT (id) DO NOTHING;
