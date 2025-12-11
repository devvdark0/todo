ALTER TABLE tasks
ADD COLUMN user_id CHAR(36);

ALTER TABLE tasks
ADD CONSTRAINT fk_tasks_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;