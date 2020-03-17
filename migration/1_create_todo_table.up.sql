CREATE TABLE IF NOT EXISTS `todo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(256) COLLATE utf8_unicode_ci DEFAULT NULL,
  `isSuccess` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO todolist.todo (id,title,description,isSuccess) VALUES 
(1,'test1','test desc',0);