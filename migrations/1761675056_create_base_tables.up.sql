CREATE TABLE IF NOT EXISTS postings(
    id INTEGER PRIMARY KEY,
	link TEXT UNIQUE NOT NULL,
	descrip TEXT NOT NULL,
	postedDate TEXT NOT NULL,
	companyName TEXT NOT NULL,
	reviewOutput TEXT NOT NULL,
	ingestDate TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS qualifications(
	id INTEGER PRIMARY KEY,
	postings_id INTEGER NOT NULL,
	qualification_text TEXT NOT NULL,
	FOREIGN KEY(postings_id) REFERENCES postings(id)
);
