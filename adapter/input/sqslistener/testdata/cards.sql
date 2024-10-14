CREATE TABLE cards.cards (
	id BIGINT UNSIGNED auto_increment NOT NULL,
	card_holder_name varchar(100) NULL,
	card_number varchar(100) NULL,
	cvv varchar(100) NULL,
	`expiry_date` varchar(100) NULL,
	created_at DATETIME NULL,
	updated_at DATETIME NULL,
	CONSTRAINT cards_pk PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;
