CREATE TABLE IF NOT EXISTS hydra_oauth2_device_code 
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(255) NOT NULL DEFAULT '',
    requested_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    client_id          VARCHAR(255) NOT NULL DEFAULT '',
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             BOOL         NOT NULL DEFAULT true,
    requested_audience TEXT         NOT NULL,
    granted_audience   TEXT         NOT NULL,
    challenge_id       VARCHAR(40)  NULL REFERENCES  hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_device_code_request_id_idx ON hydra_oauth2_device_code (request_id);
CREATE INDEX hydra_oauth2_device_code_client_id_idx ON hydra_oauth2_device_code (client_id);
CREATE INDEX hydra_oauth2_device_code_challenge_id_idx ON hydra_oauth2_device_code (challenge_id);
ALTER TABLE hydra_oauth2_device_code ADD CONSTRAINT hydra_oauth2_device_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_device_code ADD CONSTRAINT hydra_oauth2_device_code_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS hydra_oauth2_user_code 
(
    signature          VARCHAR(255) NOT NULL PRIMARY KEY,
    request_id         VARCHAR(40)  NOT NULL DEFAULT '',
    requested_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    client_id          VARCHAR(255) NOT NULL DEFAULT '',
    scope              TEXT         NOT NULL,
    granted_scope      TEXT         NOT NULL,
    form_data          TEXT         NOT NULL,
    session_data       TEXT         NOT NULL,
    subject            VARCHAR(255) NOT NULL DEFAULT '',
    active             BOOL         NOT NULL DEFAULT true,
    requested_audience TEXT         NOT NULL,
    granted_audience   TEXT         NOT NULL,
    challenge_id       VARCHAR(40)  NULL REFERENCES hydra_oauth2_consent_request_handled (challenge) ON DELETE CASCADE
);
CREATE INDEX hydra_oauth2_user_code_request_id_idx ON hydra_oauth2_user_code (request_id);
CREATE INDEX hydra_oauth2_user_code_client_id_idx ON hydra_oauth2_user_code (client_id);
CREATE INDEX hydra_oauth2_user_code_challenge_id_idx ON hydra_oauth2_user_code (challenge_id);
ALTER TABLE hydra_oauth2_user_code ADD CONSTRAINT hydra_oauth2_user_code_challenge_id_fk FOREIGN KEY (challenge_id) REFERENCES hydra_oauth2_flow(consent_challenge_id) ON DELETE CASCADE;
ALTER TABLE hydra_oauth2_user_code ADD CONSTRAINT hydra_oauth2_user_code_client_id_fk FOREIGN KEY (client_id) REFERENCES hydra_client(id) ON DELETE CASCADE;
