-- phpMyAdmin SQL Dump
-- version 4.8.5
-- https://www.phpmyadmin.net/
--
-- Хост: localhost
-- Время создания: Июн 12 2019 г., 16:04
-- Версия сервера: 10.1.38-MariaDB-0+deb9u1
-- Версия PHP: 5.6.40-8+0~20190531120521.15+stretch~1.gbpa77d1d

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `ps`
--

-- --------------------------------------------------------

--
-- Структура таблицы `modules`
--

CREATE TABLE `modules` (
  `id` int(10) UNSIGNED NOT NULL,
  `id_old` int(10) UNSIGNED DEFAULT NULL,
  `id_new` int(10) UNSIGNED DEFAULT NULL,
  `pathname_old` varchar(64) NOT NULL DEFAULT '',
  `pathname_new` varchar(64) DEFAULT '',
  `name_old` varchar(128) DEFAULT '',
  `name_new` varchar(128) DEFAULT '',
  `author_old` varchar(45) DEFAULT '',
  `author_new` varchar(45) DEFAULT '',
  `version_old` varchar(8) DEFAULT '',
  `version_new` varchar(8) DEFAULT '',
  `active_old` tinyint(1) UNSIGNED DEFAULT '0',
  `active_new` tinyint(1) UNSIGNED DEFAULT '0',
  `is_configurable_old` tinyint(1) UNSIGNED DEFAULT '0',
  `is_configurable_new` tinyint(1) UNSIGNED DEFAULT '0',
  `available_url` varchar(256) DEFAULT '',
  `description_old` varchar(256) DEFAULT '',
  `description_new` varchar(256) DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `modules`
--
ALTER TABLE `modules`
  ADD PRIMARY KEY (`id`),
  ADD KEY `name` (`pathname_old`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `modules`
--
ALTER TABLE `modules`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
