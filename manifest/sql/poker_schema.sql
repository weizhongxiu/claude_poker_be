
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

DROP DATABASE IF EXISTS `app`;
CREATE DATABASE `app` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `app`;

CREATE TABLE `users` (
  `id`           BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `uid`          VARCHAR(20)      NOT NULL                 COMMENT 'User unique number',
  `nickname`     VARCHAR(50)      NOT NULL,
  `avatar`       VARCHAR(255)     NOT NULL DEFAULT '',
  `phone`        VARCHAR(20)      NOT NULL DEFAULT '',
  `password`     VARCHAR(64)      NOT NULL DEFAULT '',
  `gender`       TINYINT          NOT NULL DEFAULT 0       COMMENT '0=unknown 1=male 2=female',
  `status`       TINYINT          NOT NULL DEFAULT 1       COMMENT '1=normal 2=banned',
  `last_login_at` DATETIME        DEFAULT NULL,
  `created_at`   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`   DATETIME         DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid`   (`uid`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User';

CREATE TABLE `user_wallets` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `chips`        BIGINT          NOT NULL DEFAULT 0        COMMENT 'Chips balance',
  `gold`         BIGINT          NOT NULL DEFAULT 0        COMMENT 'Gold balance',
  `diamonds`     INT             NOT NULL DEFAULT 0        COMMENT 'Diamond balance',
  `frozen_chips` BIGINT          NOT NULL DEFAULT 0        COMMENT 'Chips frozen on table',
  `version`      INT             NOT NULL DEFAULT 0        COMMENT 'Optimistic lock version',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User wallet';

CREATE TABLE `chip_transactions` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `type`       TINYINT         NOT NULL COMMENT '1=recharge 2=withdraw 3=win 4=lose 5=buyin 6=cashout',
  `amount`     BIGINT          NOT NULL COMMENT 'Change amount (positive=add negative=sub)',
  `balance`    BIGINT          NOT NULL COMMENT 'Balance after change',
  `ref_id`     BIGINT UNSIGNED DEFAULT NULL COMMENT 'Related id (game_id etc)',
  `remark`     VARCHAR(200)    NOT NULL DEFAULT '',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id`    (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Chip transaction log';

CREATE TABLE `user_stats` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`       BIGINT UNSIGNED NOT NULL,
  `game_type`     TINYINT         NOT NULL DEFAULT 0       COMMENT '0=all 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese',
  `stat_type`     TINYINT         NOT NULL DEFAULT 1       COMMENT '1=day 2=week 3=month 4=custom',
  `stat_date`     DATE            DEFAULT NULL             COMMENT 'Stat start date, NULL=all-time',
  `stat_end`      DATE            DEFAULT NULL             COMMENT 'Stat end date (stat_type=4 only)',
  `total_games`   INT             NOT NULL DEFAULT 0       COMMENT 'Total hands',
  `total_hands`   INT             NOT NULL DEFAULT 0       COMMENT 'Total hands participated',
  `total_sessions` INT            NOT NULL DEFAULT 0       COMMENT 'Total sessions',
  `total_wins`    INT             NOT NULL DEFAULT 0       COMMENT 'Winning hands',
  `total_profit`  BIGINT          NOT NULL DEFAULT 0       COMMENT 'Total profit/loss',
  `total_buyin`   BIGINT          NOT NULL DEFAULT 0       COMMENT 'Total buy-in',
  `total_flow`    BIGINT          NOT NULL DEFAULT 0       COMMENT 'Total flow',
  `biggest_pot`   BIGINT          NOT NULL DEFAULT 0       COMMENT 'Biggest pot won',
  `vpip`          DECIMAL(5,2)    NOT NULL DEFAULT 0.00    COMMENT 'VPIP%',
  `pfr`           DECIMAL(5,2)    NOT NULL DEFAULT 0.00    COMMENT 'PFR%',
  `wtsd`          DECIMAL(5,2)    NOT NULL DEFAULT 0.00    COMMENT 'Went to showdown%',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_type_date` (`user_id`,`game_type`,`stat_type`,`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User statistics (multi-dimensional)';

CREATE TABLE `clubs` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `club_no`      VARCHAR(20)     NOT NULL                  COMMENT 'Club number',
  `name`         VARCHAR(100)    NOT NULL,
  `logo`         VARCHAR(255)    NOT NULL DEFAULT '',
  `owner_id`     BIGINT UNSIGNED NOT NULL                  COMMENT 'Owner user_id',
  `announcement` TEXT            COMMENT 'Announcement',
  `member_count` INT             NOT NULL DEFAULT 0,
  `max_members`  INT             NOT NULL DEFAULT 200,
  `status`       TINYINT         NOT NULL DEFAULT 1        COMMENT '1=active 2=dissolved',
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`   DATETIME        DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_club_no` (`club_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Club';

CREATE TABLE `club_members` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `club_id`    BIGINT UNSIGNED NOT NULL,
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `role`       TINYINT         NOT NULL DEFAULT 3          COMMENT '1=owner 2=admin 3=member',
  `chips`      BIGINT          NOT NULL DEFAULT 0          COMMENT 'Club chips',
  `status`     TINYINT         NOT NULL DEFAULT 1          COMMENT '1=active 2=banned',
  `joined_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_club_user` (`club_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Club member';

CREATE TABLE `tables` (
  `id`                  BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT,
  `table_no`            VARCHAR(20)      NOT NULL                 COMMENT 'Table number',
  `club_id`             BIGINT UNSIGNED  DEFAULT NULL             COMMENT 'Club id, NULL=public',
  `name`                VARCHAR(100)     NOT NULL DEFAULT '',
  `has_password`        TINYINT          NOT NULL DEFAULT 0       COMMENT 'Password protected',
  `password`            VARCHAR(20)      NOT NULL DEFAULT '',
  `game_type`           TINYINT          NOT NULL DEFAULT 1       COMMENT '1=holdem 2=shortdeck 3=plo 4=sng 5=chinese',
  `blind_type`          TINYINT          NOT NULL DEFAULT 1       COMMENT '1=fixed 2=increasing',
  `small_blind`         BIGINT           NOT NULL,
  `big_blind`           BIGINT           NOT NULL,
  `ante`                BIGINT           NOT NULL DEFAULT 0,
  `straddle_enabled`    TINYINT          NOT NULL DEFAULT 0       COMMENT 'Allow straddle (2x BB)',
  `min_buyin`           BIGINT           NOT NULL,
  `max_buyin`           BIGINT           NOT NULL,
  `max_buyin_total`     BIGINT           NOT NULL DEFAULT 0       COMMENT 'Cumulative max buy-in, 0=unlimited',
  `duration`            DECIMAL(4,1)     NOT NULL DEFAULT 2.0     COMMENT 'Session duration in hours',
  `run_twice`           TINYINT          NOT NULL DEFAULT 0       COMMENT 'Feature switch: allow run-twice',
  `low_water_insurance` TINYINT          NOT NULL DEFAULT 0       COMMENT 'Low water insurance',
  `crit_gameplay`       TINYINT          NOT NULL DEFAULT 0       COMMENT 'Critical hit gameplay',
  `activity_points`     TINYINT          NOT NULL DEFAULT 0       COMMENT 'Activity points enabled',
  `auto_rebuy`          TINYINT          NOT NULL DEFAULT 0       COMMENT 'Auto rebuy/cashout',
  `buyin_approval`      TINYINT          NOT NULL DEFAULT 0       COMMENT 'Rebuy needs admin approval',
  `delay_show_card`     TINYINT          NOT NULL DEFAULT 0       COMMENT 'Delay show card',
  `random_seat`         TINYINT          NOT NULL DEFAULT 0       COMMENT 'Random seat assignment',
  `spectator_mute`      TINYINT          NOT NULL DEFAULT 0       COMMENT 'Mute spectators',
  `gps_ip_restrict`     TINYINT          NOT NULL DEFAULT 0       COMMENT 'GPS and IP restriction',
  `full_table_start`    TINYINT          NOT NULL DEFAULT 0       COMMENT 'Start only when full',
  `max_seats`           TINYINT          NOT NULL DEFAULT 9       COMMENT 'Max seats 2-10 (standard 6-9)',
  `current_players`     TINYINT          NOT NULL DEFAULT 0,
  `tag`                 VARCHAR(50)      NOT NULL DEFAULT ''      COMMENT 'Custom tag',
  `creator_id`          BIGINT UNSIGNED  NOT NULL DEFAULT 0,
  `status`              TINYINT          NOT NULL DEFAULT 1       COMMENT '1=waiting 2=playing 3=closed',
  `created_at`          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`          DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `ended_at`            DATETIME         DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_table_no` (`table_no`),
  KEY `idx_club_id`     (`club_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Game table configuration';

CREATE TABLE `table_seats` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `table_id`      BIGINT UNSIGNED NOT NULL,
  `user_id`       BIGINT UNSIGNED NOT NULL,
  `seat_no`       TINYINT         NOT NULL                COMMENT 'Seat number 1-10',
  `chips`         BIGINT          NOT NULL                COMMENT 'Chips on table',
  `status`        TINYINT         NOT NULL DEFAULT 1      COMMENT '1=seated 2=left',
  `is_sitting_out` TINYINT        NOT NULL DEFAULT 0      COMMENT 'Sitting out flag',
  `joined_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_table_seat` (`table_id`,`seat_no`),
  UNIQUE KEY `uk_table_user` (`table_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Table seat';

CREATE TABLE `room_sessions` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `table_id`        BIGINT UNSIGNED NOT NULL,
  `session_no`      VARCHAR(30)     NOT NULL                COMMENT 'Session number (unique)',
  `creator_id`      BIGINT UNSIGNED NOT NULL,
  `game_type`       TINYINT         NOT NULL DEFAULT 1      COMMENT '1=holdem 2=shortdeck 3=plo 4=sng 5=chinese',
  `small_blind`     BIGINT          NOT NULL,
  `big_blind`       BIGINT          NOT NULL,
  `total_hands`     INT             NOT NULL DEFAULT 0      COMMENT 'Total hands played',
  `total_flow`      BIGINT          NOT NULL DEFAULT 0      COMMENT 'Sum of all pots',
  `total_buyin`     BIGINT          NOT NULL DEFAULT 0      COMMENT 'Total buy-in of all players',
  `max_pot`         BIGINT          NOT NULL DEFAULT 0      COMMENT 'Biggest pot in session',
  `avg_pot`         BIGINT          NOT NULL DEFAULT 0      COMMENT 'Average pot size',
  `player_count`    TINYINT         NOT NULL DEFAULT 0,
  `spectator_count` INT             NOT NULL DEFAULT 0,
  `duration`        DECIMAL(4,1)    NOT NULL DEFAULT 0      COMMENT 'Actual duration in hours',
  `status`          TINYINT         NOT NULL DEFAULT 1      COMMENT '1=running 2=ended',
  `end_reason`      TINYINT         NOT NULL DEFAULT 0      COMMENT '0=running 1=timeout 2=manual 3=all_left',
  `started_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ended_at`        DATETIME        DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_no` (`session_no`),
  KEY `idx_table_id`   (`table_id`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_started_at` (`started_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Room session (one complete game)';

CREATE TABLE `session_players` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id`   BIGINT UNSIGNED NOT NULL,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `seat_no`      TINYINT         NOT NULL,
  `total_hands`  INT             NOT NULL DEFAULT 0     COMMENT 'Hands participated',
  `total_buyin`  BIGINT          NOT NULL DEFAULT 0     COMMENT 'Cumulative buy-in',
  `chips_final`  BIGINT          NOT NULL DEFAULT 0     COMMENT 'Final chips when session ended',
  `result`       BIGINT          NOT NULL DEFAULT 0     COMMENT 'Profit/loss = chips_final - total_buyin',
  `vpip`         DECIMAL(5,2)    NOT NULL DEFAULT 0.00  COMMENT 'VPIP% in this session',
  `win_rate`     DECIMAL(5,2)    NOT NULL DEFAULT 0.00  COMMENT 'Win rate% in this session',
  `activity_pts` INT             NOT NULL DEFAULT 0     COMMENT 'Activity points earned',
  `is_mvp`       TINYINT         NOT NULL DEFAULT 0,
  `rank`         TINYINT         NOT NULL DEFAULT 0     COMMENT 'Rank by profit desc',
  `joined_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `left_at`      DATETIME        DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_user` (`session_id`,`user_id`),
  KEY `idx_user_joined` (`user_id`,`joined_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Session player summary';

CREATE TABLE `buyin_records` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id`  BIGINT UNSIGNED NOT NULL,
  `user_id`     BIGINT UNSIGNED NOT NULL,
  `amount`      BIGINT          NOT NULL             COMMENT 'Buy-in amount',
  `type`        TINYINT         NOT NULL DEFAULT 1   COMMENT '1=initial 2=rebuy 3=cashout',
  `status`      TINYINT         NOT NULL DEFAULT 1   COMMENT '1=pending 2=approved 3=rejected',
  `approved_by` BIGINT UNSIGNED DEFAULT NULL         COMMENT 'Admin user_id who approved',
  `approved_at` DATETIME        DEFAULT NULL,
  `remark`      VARCHAR(100)    NOT NULL DEFAULT '',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_session_user` (`session_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Buy-in / rebuy / cashout record';

CREATE TABLE `table_observers` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id` BIGINT UNSIGNED NOT NULL,
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `status`     TINYINT         NOT NULL DEFAULT 1    COMMENT '1=watching 2=left',
  `joined_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `left_at`    DATETIME        DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_user` (`session_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Table observer';

CREATE TABLE `table_messages` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `session_id` BIGINT UNSIGNED NOT NULL,
  `user_id`    BIGINT UNSIGNED NOT NULL               COMMENT '0=system',
  `type`       TINYINT         NOT NULL DEFAULT 1     COMMENT '1=text 2=emoji 3=system 4=invite',
  `content`    VARCHAR(500)    NOT NULL DEFAULT '',
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='In-table chat message';

CREATE TABLE `table_invitations` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `table_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0    COMMENT 'Table id (before session)',
  `session_id` BIGINT UNSIGNED NOT NULL DEFAULT 0    COMMENT 'Session id (after start), 0=pre-session',
  `inviter_id` BIGINT UNSIGNED NOT NULL,
  `invitee_id` BIGINT UNSIGNED NOT NULL,
  `status`     TINYINT         NOT NULL DEFAULT 1    COMMENT '1=pending 2=accepted 3=rejected 4=expired',
  `expired_at` DATETIME        NOT NULL,
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_invitee_id` (`invitee_id`),
  KEY `idx_session_id` (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Table invitation';

CREATE TABLE `games` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `table_id`        BIGINT UNSIGNED NOT NULL,
  `session_id`      BIGINT UNSIGNED NOT NULL DEFAULT 0   COMMENT 'Parent session id',
  `hand_no`         VARCHAR(30)     NOT NULL              COMMENT 'Hand number (unique)',
  `shuffle_seed`    VARCHAR(64)     NOT NULL DEFAULT ''   COMMENT 'Shuffle seed, published after hand ends',
  `dealer_seat`     TINYINT         NOT NULL              COMMENT 'Dealer button seat',
  `small_blind`     BIGINT          NOT NULL,
  `big_blind`       BIGINT          NOT NULL,
  `ante`            BIGINT          NOT NULL DEFAULT 0,
  `pot`             BIGINT          NOT NULL DEFAULT 0    COMMENT 'Total pot',
  `community_cards` VARCHAR(20)     NOT NULL DEFAULT ''   COMMENT 'Board cards e.g. Ah Kd Qc Jh Ts',
  `is_split_pot`    TINYINT         NOT NULL DEFAULT 0    COMMENT 'Split pot flag',
  `run_twice_used`  TINYINT         NOT NULL DEFAULT 0    COMMENT 'Run-twice actually executed',
  `run_twice_board2` VARCHAR(10)    NOT NULL DEFAULT ''   COMMENT 'Second board for run-twice (Turn+River)',
  `stage`           TINYINT         NOT NULL DEFAULT 0    COMMENT '0=blinds 1=preflop 2=flop 3=turn 4=river 5=showdown',
  `status`          TINYINT         NOT NULL DEFAULT 1    COMMENT '1=running 2=ended',
  `started_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ended_at`        DATETIME        DEFAULT NULL,
  `duration_ms`     INT             NOT NULL DEFAULT 0    COMMENT 'Hand duration in milliseconds',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hand_no`   (`hand_no`),
  KEY `idx_session_id`      (`session_id`),
  KEY `idx_table_id`        (`table_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hand (one deal)';

CREATE TABLE `game_players` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id`        BIGINT UNSIGNED NOT NULL,
  `user_id`        BIGINT UNSIGNED NOT NULL,
  `seat_no`        TINYINT         NOT NULL,
  `position`       TINYINT         NOT NULL DEFAULT 0    COMMENT '0=BTN 1=SB 2=BB 3=UTG 4=UTG+1 5=MP 6=HJ 7=CO',
  `hole_cards`     VARCHAR(10)     NOT NULL DEFAULT ''   COMMENT 'Hole cards e.g. Ah Kd',
  `forced_bet`     BIGINT          NOT NULL DEFAULT 0    COMMENT 'Forced bet (blind+ante), separate from voluntary',
  `chips_start`    BIGINT          NOT NULL              COMMENT 'Chips at hand start',
  `chips_end`      BIGINT          DEFAULT NULL          COMMENT 'Chips at hand end',
  `total_bet`      BIGINT          NOT NULL DEFAULT 0    COMMENT 'Total voluntary bet this hand',
  `result`         BIGINT          DEFAULT NULL          COMMENT 'Profit/loss this hand',
  `best_hand`      VARCHAR(30)     NOT NULL DEFAULT ''   COMMENT 'Best 5-card combo description',
  `hand_rank`      TINYINT         NOT NULL DEFAULT 0    COMMENT 'Hand strength 1(HighCard,weakest)-10(RoyalFlush,strongest). NOTE: opposite to rules.md display rank',
  `hand_rank_desc` VARCHAR(20)     NOT NULL DEFAULT ''   COMMENT 'Hand description: High Card / One Pair / ... / Royal Flush',
  `is_winner`      TINYINT         NOT NULL DEFAULT 0,
  `fold_stage`     TINYINT         DEFAULT NULL          COMMENT 'Stage when folded (0-4)',
  `is_vpip`        TINYINT         NOT NULL DEFAULT 0    COMMENT 'Voluntarily put chips in pot preflop',
  `is_pfr`         TINYINT         NOT NULL DEFAULT 0    COMMENT 'Preflop raise flag',
  `went_to_sd`     TINYINT         NOT NULL DEFAULT 0    COMMENT 'Went to showdown flag',
  `is_show_card`   TINYINT         NOT NULL DEFAULT 0    COMMENT '0=muck 1=show cards at showdown',
  `status`         TINYINT         NOT NULL DEFAULT 1    COMMENT '1=active 2=folded 3=allin',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_game_user` (`game_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hand player';

CREATE TABLE `game_actions` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id`    BIGINT UNSIGNED NOT NULL,
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `seat_no`    TINYINT         NOT NULL,
  `stage`      TINYINT         NOT NULL              COMMENT '0=blinds 1=preflop 2=flop 3=turn 4=river',
  `action`     TINYINT         NOT NULL              COMMENT '1=fold 2=check 3=call 4=raise 5=allin 6=bet 7=blind_post 8=ante_post 9=straddle',
  `amount`     BIGINT          NOT NULL DEFAULT 0,
  `pot_after`  BIGINT          NOT NULL DEFAULT 0    COMMENT 'Pot size after this action',
  `action_seq` SMALLINT        NOT NULL DEFAULT 0    COMMENT 'Action sequence in this hand',
  `action_ms`  INT             NOT NULL DEFAULT 0    COMMENT 'Decision time in milliseconds',
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_game_id` (`game_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hand action log';

CREATE TABLE `pot_distributions` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id`      BIGINT UNSIGNED NOT NULL,
  `pot_type`     TINYINT         NOT NULL DEFAULT 1   COMMENT '1=main 2=side',
  `pot_index`    TINYINT         NOT NULL DEFAULT 0   COMMENT 'Side pot index',
  `amount`       BIGINT          NOT NULL             COMMENT 'Total pot amount',
  `winner_ids`   VARCHAR(200)    NOT NULL DEFAULT ''  COMMENT 'Winner user_ids (display only), see pot_winner_details for exact amounts',
  `winner_count` TINYINT         NOT NULL DEFAULT 1   COMMENT 'Number of winners (>1 = split pot)',
  `win_reason`   VARCHAR(50)     NOT NULL DEFAULT ''  COMMENT 'Winning hand description (display)',
  `win_rank`     TINYINT         NOT NULL DEFAULT 0   COMMENT 'Winning hand strength value (1-10)',
  PRIMARY KEY (`id`),
  KEY `idx_game_id` (`game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Pot distribution summary';

CREATE TABLE `pot_winner_details` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `distribution_id` BIGINT UNSIGNED NOT NULL              COMMENT 'pot_distributions.id',
  `game_id`         BIGINT UNSIGNED NOT NULL              COMMENT 'Redundant for query',
  `user_id`         BIGINT UNSIGNED NOT NULL,
  `amount`          BIGINT          NOT NULL              COMMENT 'Exact amount this winner receives',
  `is_split`        TINYINT         NOT NULL DEFAULT 0    COMMENT 'Part of a split pot',
  PRIMARY KEY (`id`),
  KEY `idx_distribution_id` (`distribution_id`),
  KEY `idx_game_id`         (`game_id`),
  KEY `idx_user_id`         (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Pot winner detail (one row per winner per pot)';

CREATE TABLE `hand_favorites` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`    BIGINT UNSIGNED NOT NULL,
  `game_id`    BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_game` (`user_id`,`game_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Favorite hand';

CREATE TABLE `hand_replays` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `game_id`          BIGINT UNSIGNED NOT NULL,
  `hand_index`       SMALLINT        NOT NULL DEFAULT 0   COMMENT 'Hand index in session (for replay progress)',
  `stage`            TINYINT         NOT NULL             COMMENT '1=preflop 2=flop 3=turn 4=river 5=showdown',
  `community_cards`  VARCHAR(20)     NOT NULL DEFAULT ''  COMMENT 'Board at this stage',
  `pot`              BIGINT          NOT NULL DEFAULT 0,
  `players_state`    JSON            NOT NULL             COMMENT 'Player state snapshot [{seat,chips,bet,status,hole_cards}]',
  `action_seq_start` SMALLINT        NOT NULL DEFAULT 0,
  `action_seq_end`   SMALLINT        NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_game_stage` (`game_id`,`stage`),
  KEY `idx_game_id` (`game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Hand replay stage snapshot';

CREATE TABLE `tournaments` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `club_id`         BIGINT UNSIGNED DEFAULT NULL,
  `name`            VARCHAR(100)    NOT NULL,
  `type`            TINYINT         NOT NULL DEFAULT 1    COMMENT '1=MTT 2=SNG 3=Bounty',
  `buyin`           BIGINT          NOT NULL,
  `fee`             BIGINT          NOT NULL DEFAULT 0,
  `starting_chips`  BIGINT          NOT NULL,
  `max_players`     INT             NOT NULL DEFAULT 0    COMMENT '0=unlimited',
  `current_players` INT             NOT NULL DEFAULT 0,
  `prize_pool`      BIGINT          NOT NULL DEFAULT 0,
  `status`          TINYINT         NOT NULL DEFAULT 1    COMMENT '1=registering 2=running 3=ended',
  `register_start`  DATETIME        DEFAULT NULL,
  `register_end`    DATETIME        DEFAULT NULL,
  `started_at`      DATETIME        DEFAULT NULL,
  `ended_at`        DATETIME        DEFAULT NULL,
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_club_id` (`club_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Tournament';
