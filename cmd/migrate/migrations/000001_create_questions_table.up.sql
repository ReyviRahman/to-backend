CREATE TABLE IF NOT EXISTS questions (
  id BIGSERIAL PRIMARY KEY,
  category VARCHAR(10) NOT NULL,
  question_text TEXT NOT NULL,
  question_image_url TEXT,
  options JSONB NOT NULL,
  explanation TEXT,
  explanation_image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);