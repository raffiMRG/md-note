ALTER TABLE users
  ADD COLUMN role ENUM('admin', 'penulis') NOT NULL DEFAULT 'penulis';
