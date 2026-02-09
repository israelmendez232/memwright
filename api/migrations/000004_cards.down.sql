DROP TRIGGER IF EXISTS set_cards_updated_at ON cards;
DROP INDEX IF EXISTS idx_cards_deck_id;
DROP TABLE IF EXISTS cards;
DROP TYPE IF EXISTS card_type;
