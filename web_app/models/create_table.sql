create table `user` (
    `id` bigint(20) not null auto_increment,
    `user_id` bigint(20) not null,
    `username` varchar(64) collate utf8mb4_general_ci not null,
    `password` varchar(64) collate utf8mb4_general_ci not null,
    `email` varchar(64) collate utf8mb4_general_ci,
    `gender` tinyint(4) not null default '0',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update
                current_timestamp,
    primary key (`id`),
    unique key `idx_username` (`username`) using btree,
    unique key `idx_user_id` (`user_id`) using btree
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci


create table community (
    `id` int(11) not null auto_increment,
    `community_id` int(10) unsigned not null,
    `community_name` varchar(128) collate utf8mb4_general_ci not null,
    `introduction` varchar(256) collate utf8mb4_general_ci not null,
    `create_time` timestamp not null default current_timestamp,
    `update_time` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_community_id` (`community_id`),
    unique key `idx_community_name` (`community_name`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;


drop table if exists `post`;
create table `post`(
    `id` bigint(20) not null auto_increment,
    `post_id` bigint(20) not null comment '帖子id',
    `title` varchar(128) collate utf8mb4_general_ci not null comment '标题',
    `content` varchar(8192) collate utf8mb4_general_ci not null comment '内容',
    `author_id` bigint(20) not null comment '作者的用户id',
    `community_id` bigint(20) not null comment '所属社区',
    `status` tinyint(4) not null default '1' comment '帖子状态',
    `create_time` timestamp null default current_timestamp comment '创建时间',
    `update_time` timestamp null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (`id`),
    unique key `idx_post_id` (`post_id`),
    key `idx_author_id` (`author_id`),
    key `idx_community_id` (`community_id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

create table `vote_post` (
    `id` bigint(20) not null auto_increment comment '主键ID',
    `post_id` bigint(20) not null comment '帖子ID',
    `user_id` bigint(20) not null comment '用户ID',
    `vote_type` tinyint(1) not null comment '投票类型,1表示赞成,-1表示反对',
    `create_time` timestamp not null default current_timestamp comment '投票时间',
    primary key (`id`),
    unique key `idx_post_user` (`post_id`, `user_id`) comment '确保每个用户对每个帖子只能投一次票',
    key `idx_post_id` (`post_id`) comment '加速按帖子查询投票记录',
    key `idx_user_id` (`user_id`) comment '加速按用户查询投票记录'
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci comment='帖子投票表';

insert into `community` values ('1','1','Go','Golang','2016-11-01 08:10:10','2016-11-01 08:10:10');
insert into `community` values ('2','2','Leetcode','刷题刷题刷题','2024-7-04 10:10:10','2024-7-04 10:10:10');
insert into `community` values ('3','3','Shadows Die Twice','弹刀弹刀','2024-08-12 12:15:10','2024-08-12 12:15:10');
insert into `community` values ('4','4','Elden Ring','翻滚翻滚','2024-10-08 09:09:05','2024-10-08 09:09:05');
