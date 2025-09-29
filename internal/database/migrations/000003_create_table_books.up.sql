CREATE TABLE books (
    id UUID,
    name VARCHAR(255) NOT NULL,
    total_page INTEGER NOT NULL,
    author_id UUID NOT NULL,
    photo_url VARCHAR(255),
    status status NOT NULL,
    date_complete DATE,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id),
    FOREIGN KEY(author_id) REFERENCES authors (id)
);