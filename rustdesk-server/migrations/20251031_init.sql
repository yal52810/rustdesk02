/* 这是我们所有的建表语句 */

/* 1. 创建 'peer' 表 */
create table if not exists peer (
    guid blob primary key not null,
    id varchar(100) not null,
    uuid blob not null,
    pk blob not null,
    created_at datetime not null default(current_timestamp),
    user blob,
    status tinyint,
    note varchar(300),
    info text not null
) without rowid;

/* 2. 创建 'peer' 表的索引 */
create unique index if not exists index_peer_id on peer (id);
create index if not exists index_peer_user on peer (user);
create index if not exists index_peer_created_at on peer (created_at);
create index if not exists index_peer_status on peer (status);

/* 3. 创建 'users' 表 */
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active'
);

/* 4. 创建 'subscriptions' 表 */
CREATE TABLE IF NOT EXISTS subscriptions (
    subscription_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expiration_datetime DATETIME NOT NULL,
    concurrent_limit INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

/* 5. 创建 'user_device_permissions' 表 */
CREATE TABLE IF NOT EXISTS user_device_permissions (
    permission_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    device_id TEXT NOT NULL, 
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

/* 6. 创建 'user_device_permissions' 表的索引 */
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_device ON user_device_permissions (user_id, device_id);