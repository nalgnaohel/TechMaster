CREATE TABLE dialog (
id BIGSERIAL PRIMARY KEY,
lang VARCHAR(2) NOT NULL, //vi: Vietnamese, en: English
content TEXT NOT NULL //Lưu toàn bộ nội dung hội thoại
);
-- Create word table
CREATE TABLE word (
id BIGSERIAL PRIMARY KEY,
lang VARCHAR(2) NOT NULL, //vi: Vietnamese, en: English
content TEXT NOT NULL, //Lưu gc
translate TEXT NOT NULL //Lưu dịch ra ting Anh
);
-- trung gian quan hệ  giữa word và dialog
CREATE TABLE word_dialog (
dialog_id BIGINT REFERENCES dialog(id) ON DELETE CASCADE,
word_id BIGINT REFERENCES word(id) ON DELETE CASCADE,
PRIMARY KEY (dialog_id, word_id)
);