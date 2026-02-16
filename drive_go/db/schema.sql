-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Hôte : 127.0.0.1:3306
-- Généré le : lun. 09 fév. 2026 à 13:25
-- Version du serveur : 9.1.0
-- Version de PHP : 8.3.14

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de données : `drive_app`
--

-- --------------------------------------------------------

--
-- Structure de la table `albums`
--

DROP TABLE IF EXISTS `albums`;
CREATE TABLE IF NOT EXISTS `albums` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `cover_photo_id` bigint UNSIGNED DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `cover_photo_id` (`cover_photo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `album_items`
--

DROP TABLE IF EXISTS `album_items`;
CREATE TABLE IF NOT EXISTS `album_items` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `album_id` bigint UNSIGNED NOT NULL,
  `photo_id` bigint UNSIGNED NOT NULL,
  `position` int UNSIGNED DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_album_photo` (`album_id`,`photo_id`),
  KEY `photo_id` (`photo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `album_share_links`
--

DROP TABLE IF EXISTS `album_share_links`;
CREATE TABLE IF NOT EXISTS `album_share_links` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `album_id` bigint UNSIGNED NOT NULL,
  `created_by` bigint UNSIGNED NOT NULL,
  `token` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `permission` enum('view','upload') COLLATE utf8mb4_unicode_ci DEFAULT 'view',
  `expires_at` timestamp NULL DEFAULT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `revoked_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `album_id` (`album_id`),
  KEY `created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `files`
--

DROP TABLE IF EXISTS `files`;
CREATE TABLE IF NOT EXISTS `files` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `folder_id` bigint UNSIGNED DEFAULT NULL,
  `original_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `storage_key` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `mime_type` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL,
  `size_bytes` bigint UNSIGNED NOT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `storage_key` (`storage_key`),
  KEY `user_id` (`user_id`),
  KEY `folder_id` (`folder_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `files`
--

INSERT INTO `files` (`id`, `user_id`, `folder_id`, `original_name`, `storage_key`, `mime_type`, `size_bytes`, `is_deleted`, `deleted_at`, `created_at`) VALUES
(6, 1, NULL, 'Design Patterns.pdf', '5f7ccadae4c98b772293ac42bf0d58bc.pdf', 'application/pdf', 263900, 0, NULL, '2026-02-05 09:54:33'),
(7, 1, NULL, 'Design Patterns.pdf', 'd87b5e3797d604cae7bc1447690ac2e4.pdf', 'application/pdf', 263900, 0, NULL, '2026-02-05 09:56:12'),
(8, 1, NULL, 'main.go', 'f4fba914ef3f98d3802f57a9bf6585d5.go', 'text/plain; charset=utf-8', 600, 0, NULL, '2026-02-05 13:13:37'),
(9, 1, NULL, 'go.mod', '82a275535d387a450bea337759e79524.mod', 'text/plain; charset=utf-8', 1805, 0, NULL, '2026-02-06 09:38:56'),
(10, 1, NULL, 'Letter-M-logo-by-Mansel-Brist(1).jpg', 'f412af47b19dfb810de850c8c9e93c3b.jpg', 'image/jpeg', 50911, 0, NULL, '2026-02-06 09:41:30'),
(11, 1, NULL, 'ff64729a-d38c-4e6b-bffa-8efb692fa077.webp', 'ff2aebc9552e503601568a549fdfbcf8.webp', 'image/webp', 252770, 0, NULL, '2026-02-06 10:02:19'),
(12, 1, NULL, 'WIN_20260206_11_18_25_Pro.jpg', 'f169245fbda9a9ebfd2b4af494bc0e1f.jpg', 'image/jpeg', 108391, 0, NULL, '2026-02-06 10:18:40'),
(13, 1, NULL, 'WIN_20250911_09_55_32_Pro.jpg', '2540daeb8384b9fa7753d3c445a22aff.jpg', 'image/jpeg', 103363, 0, NULL, '2026-02-06 14:10:55');

-- --------------------------------------------------------

--
-- Structure de la table `file_permissions`
--

DROP TABLE IF EXISTS `file_permissions`;
CREATE TABLE IF NOT EXISTS `file_permissions` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `permission` enum('view','edit') COLLATE utf8mb4_unicode_ci DEFAULT 'view',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_file_user` (`file_id`,`user_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `folders`
--

DROP TABLE IF EXISTS `folders`;
CREATE TABLE IF NOT EXISTS `folders` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `parent_id` bigint UNSIGNED DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `photos`
--

DROP TABLE IF EXISTS `photos`;
CREATE TABLE IF NOT EXISTS `photos` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `width` int UNSIGNED DEFAULT NULL,
  `height` int UNSIGNED DEFAULT NULL,
  `thumb_storage_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `thumb_width` int UNSIGNED DEFAULT NULL,
  `thumb_height` int UNSIGNED DEFAULT NULL,
  `taken_at` datetime DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_id` (`file_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Structure de la table `share_links`
--

DROP TABLE IF EXISTS `share_links`;
CREATE TABLE IF NOT EXISTS `share_links` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_id` bigint UNSIGNED NOT NULL,
  `token` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` datetime DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `fk_share_file` (`file_id`)
) ENGINE=MyISAM AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `share_links`
--

INSERT INTO `share_links` (`id`, `file_id`, `token`, `created_at`, `expires_at`, `is_active`) VALUES
(1, 6, 'hDrrO6Xo-wDWJht-WUBvFpKqhenyNgQk', '2026-02-06 10:25:35', NULL, 1),
(2, 6, 'svuSEfhukoUgh6usWCd2HMe6BIO3YM5-', '2026-02-06 10:43:49', NULL, 1),
(3, 6, 'cJdKmJT4AxDbYfx3Gy8D1Ee2IFsAK9wG', '2026-02-06 10:45:32', NULL, 1),
(4, 6, 'D7fWDq04kVeYcDwejQBZlU7sVIjtzdmY', '2026-02-06 10:53:50', NULL, 1),
(5, 10, '7MMHUxL-n-YfdeI16vY8RXIhfeO7TMfv', '2026-02-06 10:56:02', NULL, 1),
(6, 10, 'nSmh71K1NAorYH56J1V-qHgKOhRTSoLx', '2026-02-06 11:00:57', NULL, 1),
(7, 9, '5Aae7Rn7QWwdEfyTzl7se_0jWHg_ZFzu', '2026-02-06 11:01:58', NULL, 1),
(8, 11, 'ypUxY1gjvIcYrAH_tIeC095ED3alLLxE', '2026-02-06 11:02:24', NULL, 1),
(9, 11, 'oq6db-2q4tlfJO5u0903BRNt24D0PkmC', '2026-02-06 11:14:47', NULL, 0),
(10, 12, 'QmyimYb2fYQ0cU4ERqQpOTaAuGRFb7nl', '2026-02-06 11:18:46', NULL, 0),
(11, 12, 'Xeyws8ZzOKMD3r-bBMO5iOzaKgHXv8zL', '2026-02-06 11:19:15', NULL, 0),
(12, 12, 'ebHKoglF0g-RbBTJfySeQh0u53fG6zSA', '2026-02-06 11:25:38', NULL, 0),
(13, 12, 'oAoHQ9wLSzY8h5XOsNi8EgUpbIjryts8', '2026-02-06 11:39:39', NULL, 0),
(14, 12, '2Senm_fIIOX-JrsfQuqXThuFxmDdCCtI', '2026-02-06 11:40:00', NULL, 0),
(15, 13, 'QkqWODy6cpOAb_BSFAhehCeldwM4pyNr', '2026-02-06 15:11:35', NULL, 0);

-- --------------------------------------------------------

--
-- Structure de la table `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` varchar(190) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `full_name` varchar(120) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Déchargement des données de la table `users`
--

INSERT INTO `users` (`id`, `email`, `password_hash`, `full_name`, `created_at`) VALUES
(1, 'test@drive.local', 'x', 'Test User', '2026-02-05 09:43:02'),
(2, 'saidou@gmail.com', '$2a$10$i0/76T8lpxY142oANPGZ9.p/v4D3pJEklDDzGXFvbDWwHTcFM2pLK', 'saidou', '2026-02-06 15:04:17');

--
-- Contraintes pour les tables déchargées
--

--
-- Contraintes pour la table `albums`
--
ALTER TABLE `albums`
  ADD CONSTRAINT `albums_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `albums_ibfk_2` FOREIGN KEY (`cover_photo_id`) REFERENCES `photos` (`id`) ON DELETE SET NULL;

--
-- Contraintes pour la table `album_items`
--
ALTER TABLE `album_items`
  ADD CONSTRAINT `album_items_ibfk_1` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `album_items_ibfk_2` FOREIGN KEY (`photo_id`) REFERENCES `photos` (`id`) ON DELETE CASCADE;

--
-- Contraintes pour la table `album_share_links`
--
ALTER TABLE `album_share_links`
  ADD CONSTRAINT `album_share_links_ibfk_1` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `album_share_links_ibfk_2` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Contraintes pour la table `files`
--
ALTER TABLE `files`
  ADD CONSTRAINT `files_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `files_ibfk_2` FOREIGN KEY (`folder_id`) REFERENCES `folders` (`id`) ON DELETE SET NULL;

--
-- Contraintes pour la table `file_permissions`
--
ALTER TABLE `file_permissions`
  ADD CONSTRAINT `file_permissions_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `file_permissions_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- Contraintes pour la table `folders`
--
ALTER TABLE `folders`
  ADD CONSTRAINT `folders_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `folders_ibfk_2` FOREIGN KEY (`parent_id`) REFERENCES `folders` (`id`) ON DELETE CASCADE;

--
-- Contraintes pour la table `photos`
--
ALTER TABLE `photos`
  ADD CONSTRAINT `photos_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `photos_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
