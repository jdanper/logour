
CREATE TABLE IF NOT EXISTS logour.event (
    id timeuuid,
    saved_at timestamp,
    created_at timestamp,
    json_data text,
    client text,
    hostname text,
    remote_address text,
    kind text,
    message text,
    user_agent text,
    PRIMARY KEY (id, saved_at)
) WITH CLUSTERING ORDER BY (saved_at DESC);