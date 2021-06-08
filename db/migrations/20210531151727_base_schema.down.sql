ALTER TABLE `people` DROP FOREIGN KEY `team_id_fk`;
ALTER TABLE `turns` DROP FOREIGN KEY `person_id_fk`;

DROP TABLE `people`;
DROP TABLE `teams`;
DROP TABLE `turns`;

