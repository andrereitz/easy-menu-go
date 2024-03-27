# easy-menu-go

Easy menu is the same project as [Easy Menu](https://github.com/andrereitz/easy-menu), but written in Go and provided as a service API.

## Database Schema
```
CREATE TABLE users (
  id INTEGER PRIMARY KEY ASC,
  email TEXT UNIQUE NOT NULL,
  hash TEXT NOT NULL,
  business_name TEXT,
  business_url TEXT UNIQUE,
  business_color TEXT,
  business_logo TEXT
);

CREATE TABLE categories (
  id INTEGER PRIMARY KEY ASC,
  user INTEGER,
  title TEXT,
  
  FOREIGN KEY (user) REFERENCES id (users)
);

CREATE TABLE items (
  id INTEGER PRIMARY KEY ASC,
  category INTEGER,
  user INTEGER,
  media_id INTEGER,
  title TEXT,
  description TEXT,
  price REAL,

  FOREIGN KEY (user) REFERENCES id (users) 
  FOREIGN KEY (category) REFERENCES id (categories) 
  FOREIGN KEY (media_id) REFERENCES id (medias) 
)

CREATE TABLE medias (
    id INTEGER PRIMARY KEY ASC,
    url TEXT UNIQUE,
    alt TEXT,
    user INTEGER,

    FOREIGN KEY (user) REFERENCES id (users)
)
```