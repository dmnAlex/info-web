--- delete everything from db
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TYPE CheckStatus AS ENUM ('Start', 'Success', 'Failure');

CREATE TABLE Peers
(
  Nickname VARCHAR PRIMARY KEY,
  Birthday DATE NOT NULL
);

CREATE TABLE Tasks
(
  Title VARCHAR PRIMARY KEY
    CHECK (title SIMILAR TO '[A-Z]+\d+\_%' 
           AND title <> ParentTask),
  ParentTask VARCHAR NULL
    REFERENCES Tasks(Title)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  MaxXP BIGINT NOT NULL 
    CHECK (MaxXP > 0)
);

-- чтоб в таблице Tasks была только одна запись с ParentTask = null
CREATE UNIQUE INDEX ch_parent_task_null_exists
ON Tasks ((ParentTask IS NULL))
WHERE ParentTask IS NULL;

CREATE TABLE Checks
(
  ID SERIAL PRIMARY KEY,
  Peer VARCHAR NOT NULL 
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  Task VARCHAR NOT NULL 
    REFERENCES Tasks(Title) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  "Date" DATE NOT NULL 
);

CREATE TABLE P2P
(
  ID SERIAL PRIMARY KEY,
  "Check" BIGINT NOT NULL 
    REFERENCES Checks(ID) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CheckingPeer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  State CHECKSTATUS NOT NULL,
  "Time" TIME NOT NULL
);

-- чтобы был только один Start и только Success или Failure у одной проверки
CREATE UNIQUE INDEX idx_check_p2p_state_count
ON P2P ("Check", (State = 'Start'));

CREATE TABLE TransferredPoints
(
  ID SERIAL PRIMARY KEY,
  CheckingPeer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CheckedPeer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  PointsAmount BIGINT NOT NULL
    CHECK (PointsAmount > 0),
  CONSTRAINT ch_is_different_peers CHECK (CheckingPeer <> CheckedPeer)
);

CREATE UNIQUE INDEX idx_peers_pair_unique
ON TransferredPoints (CheckingPeer, CheckedPeer);

CREATE TABLE Friends
(
  ID SERIAL PRIMARY KEY,
  Peer1 VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  Peer2 VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CONSTRAINT ch_is_different_peers CHECK (Peer1 <> Peer2)
);

CREATE UNIQUE INDEX idx_friend1_friend2_unique
ON Friends (Peer1, Peer2);

CREATE TABLE Recommendations
(
  ID SERIAL PRIMARY KEY,
  Peer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  RecommendedPeer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  CONSTRAINT ch_is_different_peers CHECK (Peer <> RecommendedPeer)
);

CREATE UNIQUE INDEX idx_peer1_peer2_unique
ON Recommendations (Peer, RecommendedPeer);

CREATE TABLE TimeTracking
(
  ID SERIAL PRIMARY KEY,
  Peer VARCHAR NOT NULL
    REFERENCES Peers(Nickname) 
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  "Date" DATE NOT NULL,
  "Time" TIME NOT NULL,
  State SMALLINT NOT NULL
    CHECK (State IN (1, 2))
);

-- Функция для таблицы Verter (проверка, что добавляемая проверка успешна на P2P этапе)
CREATE FUNCTION is_p2p_success(check_id BIGINT) RETURNS boolean 
AS $$
SELECT EXISTS (SELECT "Check"
               FROM p2p 
               WHERE State = 'Success' AND "Check" = check_id);
$$ LANGUAGE SQL;

CREATE TABLE Verter
(
  ID SERIAL PRIMARY KEY,
  "Check" BIGINT NOT NULL 
    REFERENCES Checks(ID)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  State CHECKSTATUS NOT NULL,
  "Time" TIME NOT NULL,
  CONSTRAINT ch_check_is_p2p_success CHECK (is_p2p_success("Check"))
);

-- Чтобы был только один Start и только Success или Failure у одной проверки
CREATE UNIQUE INDEX idx_check_verter_state_count
ON Verter ("Check", (State = 'Start'));

CREATE TABLE XP
(
  ID SERIAL PRIMARY KEY,
  "Check" BIGINT NOT NULL
    REFERENCES Checks(ID)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
  XPAmount BIGINT NOT NULL
);


CREATE PROCEDURE 
  prc_import_from_csv (p_table VARCHAR, p_path VARCHAR, p_delimiter CHAR)
LANGUAGE plpgsql
AS $$
BEGIN
  EXECUTE format ('COPY %s ' || 'FROM %L ' || 'DELIMITER %L ' || 'CSV HEADER;',
                  p_table, p_path, p_delimiter);
END;
$$;

CREATE PROCEDURE
  prc_export_to_csv (p_table VARCHAR, p_path VARCHAR, p_delimiter CHAR)
LANGUAGE plpgsql
AS $$
BEGIN
  EXECUTE format ('COPY %s ' || 'TO %L ' || 'DELIMITER %L ' || 'CSV HEADER;',
                  p_table, p_path, p_delimiter);
END;
$$;


CALL prc_import_from_csv ('Peers', :data_path || '/data/Peers.csv', ',');
CALL prc_import_from_csv ('Tasks', :data_path || '/data/Tasks.csv', ',');
CALL prc_import_from_csv ('Checks', :data_path || '/data/Checks.csv', ',');
CALL prc_import_from_csv ('Friends', :data_path || '/data/Friends.csv', ',');
CALL prc_import_from_csv ('P2P', :data_path || '/data/P2P.csv', ',');
CALL prc_import_from_csv ('Recommendations', :data_path || '/data/Recommendations.csv', ',');
CALL prc_import_from_csv ('TimeTracking', :data_path || '/data/TimeTracking.csv', ',');
CALL prc_import_from_csv ('TransferredPoints', :data_path || '/data/TransferredPoints.csv', ',');
CALL prc_import_from_csv ('Verter', :data_path || '/data/Verter.csv', ',');
CALL prc_import_from_csv ('XP', :data_path || '/data/XP.csv', ',');

-- Import
DO $$
BEGIN
  PERFORM setval('timetracking_id_seq', (SELECT MAX(ID) FROM TimeTracking), true);
  PERFORM setval('verter_id_seq', (SELECT MAX(ID) FROM Verter), true);
  PERFORM setval('xp_id_seq', (SELECT MAX(ID) FROM XP), true);
  PERFORM setval('recommendations_id_seq', (SELECT MAX(ID) FROM Recommendations), true);
  PERFORM setval('friends_id_seq', (SELECT MAX(ID) FROM Friends), true);
  PERFORM setval('p2p_id_seq', (SELECT MAX(ID) FROM P2P), true);
  PERFORM setval('transferredpoints_id_seq', (SELECT MAX(ID) FROM TransferredPoints), true);
  PERFORM setval('checks_id_seq', (SELECT MAX(ID) FROM Checks), true);
END
$$;
