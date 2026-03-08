ALTER TABLE auction_items
ADD COLUMN auctioneer_id UUID NOT NULL;

ALTER TABLE auction_items 
ADD CONSTRAINT fk_user_auctions FOREIGN KEY (auctioneer_id) REFERENCES users (id) ON DELETE CASCADE;