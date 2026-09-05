CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title varchar(200) NOT NULL ,
    status VARCHAR(10) NOT NULL CHECK (status IN ('pending', 'done')),
    due_date DATE,
    created_at timestamptz not null default NOW()
);