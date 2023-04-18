-- name: CreateChat :exec
INSERT INTO
  chats (
    id,
    user_id,
    initial_message_id,
    status,
    token_usage,
    model,
    model_max_tokens,
    temperature,
    top_p,
    n,
    stop,
    max_tokens,
    presence_penalty,
    frequency_penalty
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
  );

-- name: FindChatByID :one
SELECT
  *
FROM
  chats
WHERE
  id = ?;

-- name: AddMessage :exec
INSERT INTO
  messages (
    id,
    chat_id,
    role,
    content,
    tokens,
    model,
    erased,
    order_msg
  )
VALUES
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
  );

-- name: FindMessagesByChatID :many
SELECT
  *
FROM
  messages
WHERE
  chat_id = ?
  AND erased = 0
ORDER BY
  order_msg ASC;

-- name: FindErasedMessagesByChatID :many
SELECT
  *
FROM
  messages
WHERE
  chat_id = ?
  AND erased = 1
ORDER BY
  order_msg ASC;