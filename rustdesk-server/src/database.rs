use async_trait::async_trait;
use hbb_common::{log, ResultType};
use sqlx::{
    sqlite::SqliteConnectOptions, ConnectOptions, Connection, Error as SqlxError, SqliteConnection,
};
use std::{ops::DerefMut, str::FromStr};
//use sqlx::postgres::PgPoolOptions;
//use sqlx::mysql::MySqlPoolOptions;

// use chrono::{NaiveDateTime, Utc};
type Pool = deadpool::managed::Pool<DbPool>;

pub struct DbPool {
    url: String,
}

#[async_trait]
impl deadpool::managed::Manager for DbPool {
    type Type = SqliteConnection;
    type Error = SqlxError;
    async fn create(&self) -> Result<SqliteConnection, SqlxError> {
        let mut opt = SqliteConnectOptions::from_str(&self.url).unwrap();
        opt.log_statements(log::LevelFilter::Debug);
        SqliteConnection::connect_with(&opt).await
    }
    async fn recycle(
        &self,
        obj: &mut SqliteConnection,
    ) -> deadpool::managed::RecycleResult<SqlxError> {
        Ok(obj.ping().await?)
    }
}

#[derive(Clone)]
pub struct Database {
    pool: Pool,
}

// --- MODIFIED ---
// We need to derive sqlx::FromRow so sqlx::query_as! can automatically
// build this struct from a query result.
#[derive(Default, sqlx::FromRow)]
pub struct Peer {
    pub guid: Vec<u8>,
    pub id: String,
    pub uuid: Vec<u8>,
    pub pk: Vec<u8>,
    pub user: Option<Vec<u8>>,
    pub info: String,
    pub status: Option<i64>,
}

// Simple user/subscription structures and stub implementations
// --- MODIFIED ---
// Also deriving FromRow for these.
#[derive(Clone, Debug, sqlx::FromRow)]
pub struct User {
    pub user_id: i64,
    pub status: String,
}

#[derive(Clone, Debug, sqlx::FromRow)]
pub struct Subscription {
    pub expiration_datetime: chrono::NaiveDateTime, // <-- 修正：改为 NaiveDateTime
    pub concurrent_limit: i64,
}

// 请复制从这里开始...
impl Database {
    pub async fn new(url: &str) -> ResultType<Database> {
        if !std::path::Path::new(url).exists() {
            std::fs::File::create(url).ok();
        }
        let n: usize = std::env::var("MAX_DATABASE_CONNECTIONS")
            .unwrap_or_else(|_| "1".to_owned())
            .parse()
            .unwrap_or(1);
        log::debug!("MAX_DATABASE_CONNECTIONS={}", n);
        let pool = Pool::new(
            DbPool {
                url: url.to_owned(),
            },
            n,
        );
        let _ = pool.get().await?; // test
        let db = Database { pool };
        db.create_tables().await?;
        Ok(db)
    }

    // --- MODIFIED --- (--- 已修改 ---)
    // 我们将一个大的 sqlx::query! 宏拆分为多个单独的 execute 调用。
    // 这对于 sqlx 的编译时检查器更友好。
    async fn create_tables(&self) -> ResultType<()> {
        let mut conn = self.pool.get().await?;

        // 1. 创建 'peer' 表 (您原来的代码)
        sqlx::query(
            "
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
            "
        )
        .execute(conn.deref_mut())
        .await?;

        // 2. 创建 'peer' 表的索引
        sqlx::query(
            "
            create unique index if not exists index_peer_id on peer (id);
            create index if not exists index_peer_user on peer (user);
            create index if not exists index_peer_created_at on peer (created_at);
            create index if not exists index_peer_status on peer (status);
            "
        )
        .execute(conn.deref_mut())
        .await?;

        // --- MODIFIED --- (--- 已修改 ---)
        // 3. 为 SaaS 表单独执行 CREATE

        // 创建 'users' 表
        sqlx::query(
            "
            CREATE TABLE IF NOT EXISTS users (
                user_id INTEGER PRIMARY KEY,
                username TEXT NOT NULL UNIQUE,
                password_hash TEXT NOT NULL,
                status TEXT NOT NULL DEFAULT 'active' -- e.g., 'active', 'banned'
            );
            "
        )
        .execute(conn.deref_mut())
        .await?;

        // 创建 'subscriptions' 表
        sqlx::query(
            "
            CREATE TABLE IF NOT EXISTS subscriptions (
                subscription_id INTEGER PRIMARY KEY,
                user_id INTEGER NOT NULL,
                expiration_datetime DATETIME NOT NULL,
                concurrent_limit INTEGER NOT NULL DEFAULT 1,
                FOREIGN KEY (user_id) REFERENCES users (user_id)
            );
            "
        )
        .execute(conn.deref_mut())
        .await?;

        // 创建 'user_device_permissions' 表
        sqlx::query(
            "
            CREATE TABLE IF NOT EXISTS user_device_permissions (
                permission_id INTEGER PRIMARY KEY,
                user_id INTEGER NOT NULL,
                device_id TEXT NOT NULL,
                FOREIGN KEY (user_id) REFERENCES users (user_id)
            );
            "
        )
        .execute(conn.deref_mut())
        .await?;

        // 创建 'user_device_permissions' 表的索引
        sqlx::query(
            "
            CREATE UNIQUE INDEX IF NOT EXISTS idx_user_device ON user_device_permissions (user_id, device_id);
            "
        )
        .execute(conn.deref_mut())
        .await?;

        Ok(())
    }

    // **** ⬇️⬇️⬇️ 这里是删除了重复代码的部分 ⬇️⬇️⬇️ ****
    // **** (旧的、错误的 create_tables 逻辑已被移除) ****

    pub async fn get_peer(&self, id: &str) -> ResultType<Option<Peer>> {
        Ok(sqlx::query_as::<_, Peer>(
            "select guid, id, uuid, pk, user, status, info from peer where id = ?"
        )
        .bind(id)
        .fetch_optional(self.pool.get().await?.deref_mut())
        .await?)
    }

    pub async fn insert_peer(
        &self,
        id: &str,
        uuid: &[u8],
        pk: &[u8],
        info: &str,
    ) -> ResultType<Vec<u8>> {
        let guid = uuid::Uuid::new_v4().as_bytes().to_vec();
        sqlx::query(
            "insert into peer(guid, id, uuid, pk, info) values(?, ?, ?, ?, ?)"
        )
        .bind(&guid)
        .bind(id)
        .bind(uuid)
        .bind(pk)
        .bind(info)
        .execute(self.pool.get().await?.deref_mut())
        .await?;
        Ok(guid)
    }

    pub async fn update_pk(
        &self,
        guid: &Vec<u8>,
        id: &str,
        pk: &[u8],
        info: &str,
    ) -> ResultType<()> {
        sqlx::query(
            "update peer set id=?, pk=?, info=? where guid=?"
        )
        .bind(id)
        .bind(pk)
        .bind(info)
        .bind(guid)
        .execute(self.pool.get().await?.deref_mut())
        .await?;
        Ok(())
    }

    // --- SaaS related methods ---

    // --- MODIFIED ---
    // Implemented real database lookup for user status.
    pub async fn get_user_status(&self, user_id: i32) -> ResultType<User> {
        // fetch_one ensures we get a user or return an error if the user_id
        // (from a valid token) somehow doesn't exist in the DB.
        let user = sqlx::query_as::<_, User>(
            "SELECT user_id, status FROM users WHERE user_id = ?"
        )
        .bind(user_id)
        .fetch_one(self.pool.get().await?.deref_mut())
        .await?;
        Ok(user)
    }

    // --- MODIFIED ---
    // Implemented real database lookup for subscription status.
    // This now uses the 'user_id' parameter, which fixes the
    // 'unused_variables' warning you were seeing.
    pub async fn get_subscription_status(&self, user_id: i32) -> ResultType<Option<Subscription>> {
        // fetch_optional is correct here, as a user might not have a subscription.
        let subscription = sqlx::query_as::<_, Subscription>(
            "SELECT expiration_datetime, concurrent_limit FROM subscriptions WHERE user_id = ?"
        )
        .bind(user_id)
        .fetch_optional(self.pool.get().await?.deref_mut())
        .await?;
        Ok(subscription)
    }

    // --- MODIFIED ---
    // Implemented real database lookup for device permission.
    // Note: The arguments are no longer prefixed with '_'
    pub async fn check_device_permission(
        &self,
        user_id: i32,
        target_device_id: &str,
    ) -> ResultType<bool> {
        // We just need to check if a record exists.
        // We use fetch_optional and check if the result is_some().
        let permission = sqlx::query(
            "SELECT 1 FROM user_device_permissions WHERE user_id = ? AND device_id = ?"
        )
        .bind(user_id)
        .bind(target_device_id)
        .fetch_optional(self.pool.get().await?.deref_mut())
        .await?;

        // If permission is Some, a row was found, so permission is granted (true).
        // If permission is None, no row was found, so permission is denied (false).
        Ok(permission.is_some())
    }
}
// ... 复制到这里结束 (这是 impl Database 的结束 '}')