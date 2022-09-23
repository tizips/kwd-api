-- +goose Up
-- +goose StatementBegin

CREATE TABLE `web_category`
(
    `id`          int unsigned     NOT NULL AUTO_INCREMENT,
    `theme`       varchar(5)       NOT NULL DEFAULT '' COMMENT '导航栏：light=明亮；dark=黑暗',
    `uri`         varchar(32)      NOT NULL DEFAULT '' COMMENT 'URI',
    `name`        varchar(32)      NOT NULL DEFAULT '' COMMENT '名称',
    `picture`     varchar(255)     NOT NULL DEFAULT '' COMMENT '图片',
    `title`       varchar(255)     NOT NULL DEFAULT '' COMMENT 'SEO 标题',
    `keyword`     varchar(255)     NOT NULL DEFAULT '' COMMENT 'SEO 关键词',
    `description` varchar(255)     NOT NULL DEFAULT '' COMMENT 'SEO 描述',
    `html`        text                      default null comment '',
    `is_enable`   tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `created_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp                 default null,
    PRIMARY KEY (`id`)
) AUTO_INCREMENT = 1000
  DEFAULT COLLATE = utf8mb4_unicode_ci COMMENT ='官网栏目表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists `web_category`;

-- +goose StatementEnd