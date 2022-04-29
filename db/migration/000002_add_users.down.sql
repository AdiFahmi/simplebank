ALTER TABLE `accounts` DROP FOREIGN KEY `owner_currency_key`;

ALTER TABLE `accounts` DROP FOREIGN KEY `accounts_owner_fkey`;

DROP TABLE IF EXISTS `users`;