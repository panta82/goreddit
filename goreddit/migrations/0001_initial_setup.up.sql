CREATE TABLE threads (
  id uuid PRIMARY KEY,
  title text NOT NULL,
  description text NOT NULL
);

CREATE TABLE posts (
  id uuid PRIMARY KEY,
  thread_id uuid NOT NULL REFERENCES threads(id),
  title text NOT NULL,
  content text NOT NULL,
  votes int NOT NULL DEFAULT 0
);

CREATE TABLE comments (
  id uuid PRIMARY KEY,
  post_id uuid NOT NULL REFERENCES posts(id),
  content text NOT NULL,
  votes int NOT NULL DEFAULT 0
);