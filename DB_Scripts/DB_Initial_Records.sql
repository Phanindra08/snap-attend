INSERT INTO snap_attend.user_profiles (role)
VALUES ('Student'), ('Professor'), ('Admin') ON CONFLICT DO NOTHING;

INSERT INTO snap_attend.addresses (address1, address2, city, state, zip_code, country)
VALUES ('123 Admin Street', '', 'Charlotte', 'NC', '28262', 'USA') RETURNING id;

INSERT INTO snap_attend.users (first_name, last_name, email, password, profile_id,
    address_id, created_at, updated_at)
VALUES ('Admin', 'User', 'admin@snapattend.com',
    '8pCab1++vOlTkCdw1PCc/A$8bPznd04UEd4ze4FiKCJjnKtDduCotTxRvKh0iJHhos', 3, 4, NOW(), NOW()); -- The password of admin user is Admin.123