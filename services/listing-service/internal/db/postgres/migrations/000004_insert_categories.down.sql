DELETE
FROM categories
WHERE id IN (
             'd290f1ee-6c54-4b01-90e6-d701748f0851', -- Электроника
             'a0ee5555-5555-5555-5555-555555555555', -- Транспорт
             'c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa' -- Одежда
    );

DELETE
FROM categories
WHERE level = 2
  AND parent_id IN (SELECT id
                    FROM categories
                    WHERE level = 1);

DELETE
FROM categories
WHERE level = 1
  AND parent_id IN (
                    'd290f1ee-6c54-4b01-90e6-d701748f0851',
                    'a0ee5555-5555-5555-5555-555555555555',
                    'c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    );

DELETE
FROM categories
WHERE id IN (
             'd290f1ee-6c54-4b01-90e6-d701748f0851',
             'a0ee5555-5555-5555-5555-555555555555',
             'c300aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'
    );
