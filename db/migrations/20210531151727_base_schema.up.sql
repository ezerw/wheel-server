CREATE TABLE `teams`
(
    `id`         bigint AUTO_INCREMENT PRIMARY KEY,
    `name`       varchar(100) UNIQUE NOT NULL,
    `created_at` timestamp default now(),
    `updated_at` timestamp default now()
);

CREATE TABLE `people`
(
    `id`         bigint AUTO_INCREMENT PRIMARY KEY,
    `first_name` varchar(100)       NOT NULL,
    `last_name`  varchar(100)       NOT NULL,
    `email`      varchar(80) UNIQUE NOT NULL,
    `team_id`    BIGINT             NOT NULL,
    `created_at` timestamp default now(),
    `updated_at` timestamp default now()
);

CREATE TABLE `turns`
(
    `id`         bigint AUTO_INCREMENT PRIMARY KEY,
    `person_id`  bigint NOT NULL,
    `date`       date   NOT NULL,
    `created_at` timestamp default now(),
    `updated_at` timestamp default now()
);

ALTER TABLE `people`
    ADD CONSTRAINT team_id_fk
        FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE;

ALTER TABLE `turns`
    ADD CONSTRAINT person_id_fk
        FOREIGN KEY (`person_id`) REFERENCES `people` (`id`) ON DELETE CASCADE;

CREATE UNIQUE INDEX `turns_index_0` ON `turns` (`date`, `person_id`);