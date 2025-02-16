CREATE TABLE dialog (
    id BIGSERIAL PRIMARY KEY,
    lang VARCHAR(2) NOT NULL, 
    content TEXT NOT NULL 
);
-- Create word table
CREATE TABLE word (
    id BIGSERIAL PRIMARY KEY,
    lang VARCHAR(2) NOT NULL, 
    content TEXT NOT NULL, 
    translate TEXT NOT NULL 
);
-- trung gian quan hệ  giữa word và dialog
CREATE TABLE word_dialog (
    dialog_id BIGINT REFERENCES dialog(id) ON DELETE CASCADE,
    word_id BIGINT REFERENCES word(id) ON DELETE CASCADE,
    PRIMARY KEY (dialog_id, word_id)
);