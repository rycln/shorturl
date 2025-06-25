package storage

const sqlAddURLPair = `
	INSERT INTO urls 
	(user_id, short_url, original_url) 
	VALUES ($1, $2, $3)
`

const sqlGetURLPairByShort = `
	SELECT 
		user_id,
		original_url, 
		is_deleted 
	FROM urls 
	WHERE short_url = $1
`

const sqlGetURLPairBatchByUserID = `
	SELECT 
		user_id, 
		short_url, 
		original_url 
	FROM urls 
	WHERE user_id = $1
`

const sqlDeleteRequestedURLs = `
	UPDATE urls 
	SET is_deleted = TRUE 
	WHERE short_url = $1
`

const sqlGetStats = `
SELECT 
	(SELECT COUNT(*) FROM urls) AS total_urls,
    (SELECT COUNT(*) FROM users) AS total_users;

`
