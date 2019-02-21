
CREATE TABLE IF NOT EXISTS logour.event (
    id timeuuid,
    saved_at timestamp,
    action text,
    created_at timestamp,
    custom_data map<text, text>,
    hostname text,
    ip text,
    kind text,
    message text,
    user_agent text,
    PRIMARY KEY (id, saved_at)
) WITH CLUSTERING ORDER BY (saved_at DESC)
