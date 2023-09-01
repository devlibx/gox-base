# How to run

1. Setup your database as given in following section
2. Set the following ENV var
   ```
   DB_NAME=<your database name>
   DB_USER=<your user>
   DB_PASS=<your password>
   DB_URL=<DB URL>
   ```
3. Run the example `queue/example/example.go`

# Database

### DB Schema

These are the tables to be created for this queue implementation:

1. jobs - this table contains all job scheduling data
2. jobs_data - this table contains the user data for the job e.g. user udf, metadata etc

Note - this table as a column `archive_after` which is default to `process_at` + 24Hr. This column
will be used to partition the data and then these partitions can be dropped after `archive_after`.

It is important that you make sure that  `process_at` never goes out of `archive_after` otherwise
it will be deleted before processing.
When can `process_at` go out of `archive_after`:

1. Yn case of error you may want to reschedule it, and you passed the next retry time > `archive_after`
2. Your system is down for some time and could not process the jobs as of now

```sql
CREATE TABLE `jobs`
(
    `id`                varchar(40)      NOT NULL,
    `tenant`            TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `correlation_id`    varchar(128)              DEFAULT NULL,
    `job_type`          TINYINT UNSIGNED NOT NULL DEFAULT '1',
    `state`             TINYINT UNSIGNED NOT NULL DEFAULT '1',
    `sub_state`         TINYINT UNSIGNED NOT NULL DEFAULT '11',
    `pending_execution` TINYINT UNSIGNED NOT NULL DEFAULT '3',
    `version`           TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `process_at`        timestamp        NOT NULL,
    `part`              timestamp        NOT NULL,
    `created_at`        timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`, `part`),
    KEY `process_at_index` (`process_at`, `job_type`, `state`, `tenant`, `pending_execution`),
    KEY `job_type_index` (`job_type`, `state`, `tenant`)
) PARTITION BY RANGE (UNIX_TIMESTAMP(`part`)) (
        PARTITION p202309_week1 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-04')), -- Week 1 (Sep 2023)
        PARTITION p202309_week2 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-11')), -- Week 2 (Sep 2023)
        PARTITION p202309_week3 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-18')), -- Week 3 (Sep 2023)
        PARTITION p202309_week4 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-25')), -- Week 4 (Sep 2023)
        PARTITION p202310_week1 VALUES LESS THAN (UNIX_TIMESTAMP('2023-10-02'))
        );


CREATE TABLE `jobs_data`
(
    `id`           varchar(40)      NOT NULL,
    `tenant`       TINYINT UNSIGNED NOT NULL DEFAULT '0',
    `properties`   text                      DEFAULT NULL,
    `string_udf_1` text,
    `string_udf_2` text,
    `int_udf_1`    int                       DEFAULT NULL,
    `int_udf_2`    int                       DEFAULT NULL,
    `part`         timestamp        NOT NULL,
    `created_at`   timestamp        NULL     DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   timestamp        NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`, `part`)
) PARTITION BY RANGE (UNIX_TIMESTAMP(`part`)) (
    PARTITION p202309_week1 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-04')), -- Week 1 (Sep 2023)
    PARTITION p202309_week2 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-11')), -- Week 2 (Sep 2023)
    PARTITION p202309_week3 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-18')), -- Week 3 (Sep 2023)
    PARTITION p202309_week4 VALUES LESS THAN (UNIX_TIMESTAMP('2023-09-25')), -- Week 4 (Sep 2023)
    PARTITION p202310_week1 VALUES LESS THAN (UNIX_TIMESTAMP('2023-10-02'))
    );
```