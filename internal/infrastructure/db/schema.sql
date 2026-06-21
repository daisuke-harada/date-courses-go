-- テーブル: users
CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  gender VARCHAR(255) NOT NULL,
  image VARCHAR(255),
  admin TINYINT(1) NOT NULL DEFAULT 0,
  password_digest VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_users_name (name),
  UNIQUE KEY uq_users_email (email)
);

-- テーブル: date_spots
CREATE TABLE date_spots (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  genre_id INT,
  prefecture_id INT,
  name VARCHAR(255) NOT NULL,
  city_name VARCHAR(255) NOT NULL,
  image VARCHAR(255),
  latitude DOUBLE,
  longitude DOUBLE,
  source VARCHAR(20) NOT NULL DEFAULT 'manual',
  maps_url VARCHAR(1000),
  normalized_name VARCHAR(255),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

-- indexes (date_spots)
CREATE INDEX index_date_spots_on_genre_id_and_created_at ON date_spots (genre_id, created_at);
CREATE INDEX index_date_spots_on_prefecture_id_and_created_at ON date_spots (prefecture_id, created_at);
CREATE INDEX index_date_spots_on_normalized_name_and_prefecture_id ON date_spots (normalized_name, prefecture_id);

-- テーブル: courses
CREATE TABLE courses (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  travel_mode VARCHAR(255) NOT NULL,
  authority VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_courses_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- indexes (courses)
CREATE INDEX index_courses_on_user_id ON courses (user_id);

-- テーブル: date_spot_reviews
CREATE TABLE date_spot_reviews (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  rate FLOAT,
  content TEXT,
  user_id BIGINT UNSIGNED NOT NULL,
  date_spot_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_date_spot_reviews_user_date_spot (user_id, date_spot_id),
  CONSTRAINT fk_date_spot_reviews_date_spots FOREIGN KEY (date_spot_id) REFERENCES date_spots (id),
  CONSTRAINT fk_date_spot_reviews_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- indexes (date_spot_reviews)
CREATE INDEX index_date_spot_reviews_on_date_spot_id ON date_spot_reviews (date_spot_id);
CREATE INDEX index_date_spot_reviews_on_user_id ON date_spot_reviews (user_id);

-- テーブル: during_spots
CREATE TABLE during_spots (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  course_id BIGINT UNSIGNED NOT NULL,
  date_spot_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_during_spots_courses FOREIGN KEY (course_id) REFERENCES courses (id),
  CONSTRAINT fk_during_spots_date_spots FOREIGN KEY (date_spot_id) REFERENCES date_spots (id)
);

-- indexes (during_spots)
CREATE INDEX index_during_spots_on_course_id ON during_spots (course_id);
CREATE INDEX index_during_spots_on_date_spot_id ON during_spots (date_spot_id);

-- テーブル: relationships
CREATE TABLE relationships (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  follow_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_relationships_users FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT fk_relationships_follow_users FOREIGN KEY (follow_id) REFERENCES users (id)
);

-- indexes (relationships)
CREATE INDEX index_relationships_on_follow_id ON relationships (follow_id);
CREATE INDEX index_relationships_on_user_id ON relationships (user_id);
