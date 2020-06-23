CREATE TABLE users (
  uuid            SERIAL,
  email           varchar(255)  NOT NULL,
  email_verified  number(1),
  premium         number(1),
  fullname        varchar(255)  NOT NULL,
  password_hash   varchar(255),
  password_hint   varchar(255),
  key_pass        varchar(255),
  private_key     binary,
  public_key      binary,
  totp_secret     varchar(255),
  security_stamp  varchar(255),
  created_at      datetime,
  PRIMARY KEY (uuid),
  UNIQUE(uuid, email)
);

CREATE TABLE devices (
  uuid        SERIAL,
  type        number(1),
  name        varchar(255),
  push_token  varchar(256),
  access_token varchar(256),
  refresh_token varchar(256),
  token_expires_at datetime,
  user_uuid   SERIAL,
  PRIMARY KEY (uuid),
  FOREIGN KEY (user_uuid) REFERENCES users(uuid)  
  UNIQUE(uuid)
)
