CREATE TABLE TweetV2 (
    ID STRING(MAX) NOT NULL,
    Author STRING(MAX) NOT NULL,
    CommitedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    Content STRING(MAX),
    Count INT64,
    CreatedAt TIMESTAMP NOT NULL,
    Favos ARRAY<STRING(MAX)>,
    ShardCreatedAt INT64,
    Sort INT64,
    UpdatedAt TIMESTAMP NOT NULL,
) PRIMARY KEY (ID)