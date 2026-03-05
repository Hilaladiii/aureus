CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) DEFAULT 'USER' CHECK (role IN ('ADMIN', 'USER')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallets(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL,
    active_balance DECIMAL(15,2) DEFAULT 0,
    held_balance DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP    
);

CREATE TABLE categories(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auction_items(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_price DECIMAL(15,2) NOT NULL,
    current_price DECIMAL(15,2) DEFAULT 0,
    start_time TIMESTAMP,
    end_time TIMESTAMP  
);

CREATE TABLE bid_histories(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL,
    user_id UUID NOT NULL,
    bid_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('WINNING', 'OUTBID', 'WON')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets 
ADD CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE auction_items
ADD CONSTRAINT fk_item_category FOREIGN KEY (category_id) REFERENCES categories (id);

ALTER TABLE bid_histories 
ADD CONSTRAINT fk_bid_item FOREIGN KEY (item_id) REFERENCES auction_items (id),
ADD CONSTRAINT fk_bid_user FOREIGN KEY (user_id) REFERENCES users (id);

