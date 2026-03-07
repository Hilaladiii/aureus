CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE Table auction_images(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP    
);

ALTER TABLE auction_images 
ADD CONSTRAINT fk_auction_images FOREIGN KEY (item_id) REFERENCES auction_items (id) ON DELETE CASCADE;
