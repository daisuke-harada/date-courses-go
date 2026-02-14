-- テーブル: users
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- テーブル: date_spots
CREATE TABLE date_spots (
  id SERIAL PRIMARY KEY,
  genre_id INTEGER,
  name VARCHAR(255) NOT NULL,
  image VARCHAR(255),
  opening_time TIMESTAMP,
  closing_time TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- indexes (date_spots)
CREATE INDEX index_date_spots_on_genre_id_and_created_at ON date_spots (genre_id, created_at);

-- テーブル: addresses
CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  prefecture_id INTEGER,
  date_spot_id INTEGER,
  city_name VARCHAR(255) NOT NULL,
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_addresses_date_spots FOREIGN KEY (date_spot_id) REFERENCES date_spots (id)
  -- NOTE: prefectures テーブルを追加したら、以下の外部キーを追加してください
  -- , CONSTRAINT fk_addresses_prefectures FOREIGN KEY (prefecture_id) REFERENCES prefectures (id)
);

-- indexes (addresses)
CREATE INDEX index_addresses_on_date_spot_id_and_created_at ON addresses (date_spot_id, created_at);
-- NOTE: prefectures テーブルを追加したら有効化してください
-- CREATE INDEX index_addresses_on_prefecture_id_and_created_at ON addresses (prefecture_id, created_at);

-- テーブル: courses
CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  travel_mode VARCHAR(255) NOT NULL,
  authority VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_courses_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- indexes (courses)
CREATE INDEX index_courses_on_user_id ON courses (user_id);

-- テーブル: date_spot_reviews
CREATE TABLE date_spot_reviews (
  id SERIAL PRIMARY KEY,
  rate FLOAT,
  content TEXT,
  user_id BIGINT NOT NULL,
  date_spot_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, date_spot_id),
  CONSTRAINT fk_date_spot_reviews_date_spots FOREIGN KEY (date_spot_id) REFERENCES date_spots (id),
  CONSTRAINT fk_date_spot_reviews_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- indexes (date_spot_reviews)
CREATE INDEX index_date_spot_reviews_on_date_spot_id ON date_spot_reviews (date_spot_id);
CREATE INDEX index_date_spot_reviews_on_user_id ON date_spot_reviews (user_id);

-- テーブル: during_spots
CREATE TABLE during_spots (
  id SERIAL PRIMARY KEY,
  course_id BIGINT NOT NULL,
  date_spot_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_during_spots_courses FOREIGN KEY (course_id) REFERENCES courses (id),
  CONSTRAINT fk_during_spots_date_spots FOREIGN KEY (date_spot_id) REFERENCES date_spots (id)
);

-- indexes (during_spots)
CREATE INDEX index_during_spots_on_course_id ON during_spots (course_id);
CREATE INDEX index_during_spots_on_date_spot_id ON during_spots (date_spot_id);

-- テーブル: relationships
CREATE TABLE relationships (
  id SERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  follow_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_relationships_users FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_relationships_follow_users FOREIGN KEY (follow_id) REFERENCES users (id)
);

-- indexes (relationships)
CREATE INDEX index_relationships_on_follow_id ON relationships (follow_id);
CREATE INDEX index_relationships_on_user_id ON relationships (user_id);
