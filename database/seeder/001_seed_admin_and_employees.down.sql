-- 001_seed_admin_and_employees.down.sql
DELETE FROM employees WHERE email LIKE 'employee%@example.com';
DELETE FROM admins WHERE email = 'admin@example.com';
