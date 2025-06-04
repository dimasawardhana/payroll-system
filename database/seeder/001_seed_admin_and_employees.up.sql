-- 001_seed_admin_and_employees.up.sql
INSERT INTO admins (name, email, password_hash, role)
VALUES ('Admin User', 'admin@example.com', '$2a$10$NT0RxZqnXHqN7bjNO3tetOCEnORymKN0SkLLvSCr0NUx6QuYrZYC.', 'admin'); -- password is 'admin123'

DO $$
DECLARE
    i INT := 1;
BEGIN
    WHILE i <= 100 LOOP
        INSERT INTO employees (name, email, password_hash, role, salary)
        VALUES (
            'Employee ' || i,
            'employee' || i || '@example.com',
            '$2a$10$zNZAJAFJmor8epLyxMdJ0eg58.4D5IucD8Y2WtpAUW7689VGnRwX2', -- password is 'employee1'
            'employee',
            5000000 + 10000 * i -- set default salary, adjust as needed
        );
        i := i + 1;
    END LOOP;
END $$;