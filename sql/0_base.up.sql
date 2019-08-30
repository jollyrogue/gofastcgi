CREATE TABLE IF NOT EXISTS form_contact (
    id INT UNSIGNED AUTO_INCREMENT,
    submission TIMESTAMP DEFAULT NOW(),
    name_first VARCHAR(128) NOT NULL,
    name_last VARCHAR(128) NOT NULL,
    phone_number VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    city VARCHAR(64) NOT NULL,
    zip_code VARCHAR(128) DEFAULT '',
    in_home_tutoring BOOLEAN DEFAULT False,
    public_location BOOLEAN DEFAULT False,
    student_grade_level ENUM('K-3', '4-6', '7-8', '9-12', 'College', 'Adult', 'Other') NOT NULL,
    tutoring_subjects VARCHAR(1024) NOT NULL,
    message TEXT DEFAULT '',
    PRIMARY KEY HASH (id)
) ENGINE=InnoDB;
