CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(20) UNIQUE NOT NULL,
  `password` varchar(256) NOT NULL,
  `secret` varchar(256) UNIQUE NOT NULL
);

CREATE TABLE `user_roles` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `role_id` int NOT NULL,
  `user_id` int NOT NULL
);

CREATE TABLE `roles` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) UNIQUE NOT NULL
);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

INSERT INTO users(`username`, `password`, `secret`) values ("admin", "admin", "admin"), ("bao", "123", "aloha");
