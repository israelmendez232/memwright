DROP TRIGGER IF EXISTS set_decks_updated_at ON decks;
DROP INDEX IF EXISTS idx_decks_parent_deck_id;
DROP INDEX IF EXISTS idx_decks_user_id;
DROP TABLE IF EXISTS decks;
