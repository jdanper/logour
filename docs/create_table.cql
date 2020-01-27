
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
    PRIMARY KEY (kind, created_at, client)
) WITH CLUSTERING ORDER BY (created_at DESC);