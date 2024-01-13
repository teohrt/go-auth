-- Users table
CREATE TABLE IF NOT EXISTS Users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(20),
    cognito_id VARCHAR(40),
    email VARCHAR(100),
    created_at TIMESTAMP default now(),
    updated_at TIMESTAMP,
    UNIQUE(username),
    UNIQUE(cognito_id),
    UNIQUE(email)
);

-- Subject_Users table
CREATE TABLE IF NOT EXISTS Subject_Users (
    subject_user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES Users(user_id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Templates table
CREATE TABLE IF NOT EXISTS Templates (
    template_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_name VARCHAR(255)
);

-- Memoirs table
CREATE TABLE IF NOT EXISTS Memoirs (
    memoir_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subject_user_id UUID REFERENCES Subject_Users(subject_user_id),
    template_id UUID REFERENCES Templates(template_id)
);

-- Template_Questions table
CREATE TABLE IF NOT EXISTS Template_Questions (
    template_id UUID REFERENCES Templates(template_id),
    question_id UUID REFERENCES Questions(question_id),
    PRIMARY KEY (template_id, question_id)
);

-- Questions table
CREATE TABLE IF NOT EXISTS Questions (
    question_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    prompt TEXT
);

-- Answers table
CREATE TABLE IF NOT EXISTS Answers (
    answer_ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID REFERENCES Questions(question_id),
    template_id UUID REFERENCES Templates(template_id),
    response TEXT
);
