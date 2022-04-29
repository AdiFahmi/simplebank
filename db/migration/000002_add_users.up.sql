CREATE TABLE `users` (
	`id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255) UNIQUE,
  `hashed_password` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE, 
  `password_changed_at` datetime NOT NULL DEFAULT '0001-01-01 00:00:01',  
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `accounts` ADD CONSTRAINT `accounts_owner_fkey` FOREIGN KEY (`owner`) REFERENCES `users` (`username`);

-- CREATE UNIQUE INDEX ON `accounts` (`owner`, `currency`);
ALTER TABLE `accounts` ADD CONSTRAINT `owner_currency_key` UNIQUE (`owner`, `currency`);
