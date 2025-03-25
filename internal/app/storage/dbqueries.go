package storage

const sqlCreateURLsTable = "CREATE TABLE IF NOT EXISTS urls (user_id UUID, short_url VARCHAR(7), original_url TEXT UNIQUE, is_deleted BOOL)"

const sqlInsertURL = "INSERT INTO urls (user_id, short_url, original_url) VALUES ($1, $2, $3)"

const sqlGetOrigURL = "SELECT original_url FROM urls WHERE short_url = $1"

const sqlGetShortURL = "SELECT short_url FROM urls WHERE original_url = $1"

const sqlGetAllUserURLs = "SELECT user_id, short_url, original_url FROM urls WHERE user_id = $1"

const sqlDeleteUserURLs = "UPDATE urls SET is_deleted = TRUE WHERE user_id = $1 AND short_url = $2"
