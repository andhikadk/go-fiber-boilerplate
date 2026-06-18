INSERT INTO users (email, password, password_is_set_by_user, role, is_active)
VALUES (
    'admin@example.com',
    '$2a$10$slYQmyNdGzin7olVN3VN2OPST9/PgBkqquzi.Ss8KIUgO2t0jWMUe',
    true,
    'admin',
    true
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (user_id, first_name, last_name)
SELECT id, 'Admin', 'User'
FROM users
WHERE email = 'admin@example.com'
ON CONFLICT (user_id) DO NOTHING;
