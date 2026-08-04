ALTER TABLE orders ADD COLUMN IF NOT EXISTS total_minor BIGINT CHECK (total_minor IS NULL OR total_minor >= 0);
ALTER TABLE orders ADD COLUMN IF NOT EXISTS currency TEXT CHECK (currency IS NULL OR char_length(currency) = 3);
