-- Users table
CREATE TABLE Users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
);

-- Subject_Users table
CREATE TABLE Subject_Users (
    subject_users_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES Users(user_id)
);

-- Templates table
CREATE TABLE Templates (
    template_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_name VARCHAR(255)
);

-- Memoirs table
CREATE TABLE Memoirs (
    memoir_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subject_users_id UUID REFERENCES Subject_Users(subject_users_id),
    template_id UUID REFERENCES Templates(template_id)
);

-- Template_Questions table
CREATE TABLE Template_Questions (
    template_id UUID REFERENCES Templates(template_id),
    question_id UUID REFERENCES Questions(question_id),
    PRIMARY KEY (template_id, question_id)
);

-- Questions table
CREATE TABLE Questions (
    question_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    prompt TEXT
);

-- Answers table
CREATE TABLE Answers (
    answer_ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID REFERENCES Questions(question_id),
    template_id UUID REFERENCES Templates(template_id),
    response TEXT
);
