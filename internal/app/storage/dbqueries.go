package storage

const sqlCreateURLsTable = "CREATE TABLE IF NOT EXISTS urls (short_url VARCHAR(7), original_url TEXT UNIQUE)"

const sqlInsertURL = "INSERT INTO urls (short_url, original_url) VALUES ($1, $2)"

const sqlGetOrigURL = "SELECT original_url FROM urls WHERE short_url = $1"

const sqlGetShortURL = "SELECT short_url FROM urls WHERE original_url = $1"
