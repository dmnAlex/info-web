/* 1) Написать процедуру добавления P2P проверки.
 Параметры: ник проверяемого, ник проверяющего, название задания, статус P2P проверки, время. 
 Если задан статус "начало", добавить запись в таблицу Checks (в качестве даты использовать сегодняшнюю). 
 Добавить запись в таблицу P2P. 
 Если задан статус "начало", в качестве проверки указать только что добавленную запись, 
 иначе указать проверку с незавершенным P2P этапом.
*/
CREATE PROCEDURE prc_add_p2p_check(p_checked_peer VARCHAR, 
                                   p_checking_peer VARCHAR,
                                   p_title VARCHAR,
                                   p_state CHECKSTATUS,
                                   p_time TIME)
AS $$
DECLARE                      
  p_check_other_id BIGINT := (
                                SELECT Checks.id
                                FROM Checks
                                INNER JOIN P2P
                                  ON Checks.id = P2P."Check"
                                  AND Checks.peer = p_checked_peer
                                  AND Checks.task = p_title
                                  AND P2P.checkingpeer = p_checking_peer
                                GROUP BY Checks.id
                                HAVING COUNT(P2P.state) = 1
                              );
  p_check_max_id BIGINT := (SELECT COALESCE(MAX(ID) + 1, 1) FROM Checks);
  p_check_start_id BIGINT := (SELECT COALESCE(p_check_other_id, p_check_max_id));
BEGIN
  IF p_state = 'Start' THEN
    INSERT INTO Checks 
    VALUES (p_check_start_id, p_checked_peer, p_title, CURRENT_DATE);
    INSERT INTO P2P
    VALUES (DEFAULT, p_check_start_id, p_checking_peer, p_state, p_time);
  ELSE
    INSERT INTO P2P VALUES (DEFAULT, p_check_other_id, p_checking_peer, p_state, p_time);
  END IF;
END;
$$ LANGUAGE plpgsql;

/* 2) Написать процедуру добавления проверки Verter'ом.
 Параметры: ник проверяемого, название задания, статус проверки Verter'ом, время. 
 Добавить запись в таблицу Verter (в качестве проверки указать проверку соответствующего 
 задания с самым поздним (по времени) успешным P2P этапом)
*/
CREATE PROCEDURE prc_add_verter_check(p_checked_peer VARCHAR, 
                                      p_title VARCHAR,
                                      p_state CHECKSTATUS,
                                      p_time TIME)
AS $$
DECLARE
last_success_check_id BIGINT :=  (
                                    SELECT Checks.ID
                                    FROM Checks 
                                    INNER JOIN P2P 
                                      ON Checks.ID = P2P."Check"
                                      AND P2P.State = 'Success'
                                      AND Checks.Task = p_title
                                      AND Checks.Peer = p_checked_peer
                                    ORDER BY "Date" DESC, "Time" DESC
                                    LIMIT 1
                                  );
BEGIN
INSERT INTO Verter VALUES (DEFAULT, last_success_check_id, p_state, p_time);
END;
$$ LANGUAGE plpgsql;

/* 3) Написать триггер: после добавления записи со статутом "начало" в таблицу P2P,
 изменить соответствующую запись в таблице TransferredPoints.
*/
CREATE FUNCTION fnc_trg_p2p_insert_transferred_points()
RETURNS TRIGGER
AS $$
DECLARE
  p_checking_peer VARCHAR := NEW.CheckingPeer;
  p_checked_peer VARCHAR := (SELECT Peer FROM Checks WHERE ID = NEW."Check");
BEGIN
  IF NEW.State = 'Start'::CHECKSTATUS
  THEN
    IF 0 = (
      SELECT COUNT(*)
      FROM TransferredPoints
      WHERE 
        CheckingPeer = p_checking_peer
        AND CheckedPeer = p_checked_peer)
    THEN
      INSERT INTO TransferredPoints
      VALUES (DEFAULT, p_checking_peer, p_checked_peer, 1);
    ELSE
      UPDATE TransferredPoints
      SET PointsAmount = PointsAmount + 1
      WHERE 
        CheckingPeer = p_checking_peer
        AND CheckedPeer = p_checked_peer;
    END IF;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_p2p_insert_transferred_points
AFTER INSERT ON P2P FOR EACH ROW
EXECUTE FUNCTION fnc_trg_p2p_insert_transferred_points();

/* 4) Написать триггер: перед добавлением записи в таблицу XP, проверить корректность 
 добавляемой записи
*/
CREATE FUNCTION fnc_trg_xp_insert()
RETURNS TRIGGER
AS $$
BEGIN
  IF NEW.XPAmount > (
    SELECT t.MaxXP
    FROM Tasks t
    INNER JOIN Checks c ON c.Task = t.Title
    WHERE c.ID = NEW."Check")
  THEN
    RAISE EXCEPTION 'XP amount bigger than Max XP for this task';
  END IF;
  IF 0 = (
    SELECT COUNT(*)
    FROM P2P
    WHERE "Check" = NEW."Check"
      AND State = 'Success'::CHECKSTATUS)
  THEN
    RAISE EXCEPTION 'p2p check is not success';
  END IF;
  IF 1 = (
    SELECT COUNT(*)
    FROM Verter
    WHERE "Check" = NEW."Check"
      AND State = 'Start'::CHECKSTATUS)
  AND 0 = (
    SELECT COUNT(*)
    FROM Verter
    WHERE "Check" = NEW."Check"
      AND State = 'Success'::CHECKSTATUS)
  THEN
    RAISE EXCEPTION 'verter check is not success';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_xp_insert
BEFORE INSERT ON XP FOR EACH ROW
EXECUTE FUNCTION fnc_trg_xp_insert();

/* 1) Написать функцию, возвращающую таблицу TransferredPoints в более человекочитаемом виде.
 Ник пира 1, ник пира 2, количество переданных пир поинтов.
 Количество отрицательное, если пир 2 получил от пира 1 больше поинтов.
*/
DROP FUNCTION IF EXISTS fnc_transferred_points_readable;
CREATE FUNCTION fnc_transferred_points_readable()
RETURNS TABLE (
  Peer1 VARCHAR,
  Peer2 VARCHAR,
  PointsAmount BIGINT
)
LANGUAGE SQL
AS $$
WITH cte AS (
  SELECT
    tp.checkingpeer AS Peer1,
    tp.checkedpeer AS Peer2,
    tp.pointsamount - COALESCE((
                                SELECT pointsamount
                                FROM transferredpoints
                                WHERE
                                  checkingpeer = tp.checkedpeer
                                  AND checkedpeer = tp.checkingpeer
                                ),
                              0) AS PointsAmount
  FROM transferredpoints tp
)
SELECT *
FROM cte
EXCEPT
SELECT Peer2, Peer1, - PointsAmount
FROM cte
WHERE
  Peer1 < Peer2
  AND Peer2 = (
                SELECT checkingpeer
                FROM transferredpoints
                WHERE
                  checkingpeer = Peer2
                  AND checkedpeer = Peer1)
$$;

/* 2) Написать функцию, которая возвращает таблицу вида:
    ник пользователя,
    название проверенного задания,
    кол-во полученного XP
 В таблицу включать только задания, успешно прошедшие проверку (определять по таблице Checks).
 Одна задача может быть успешно выполнена несколько раз. В таком случае в таблицу включать все успешные проверки.
*/
DROP FUNCTION IF EXISTS fnc_peer_task_xp;
CREATE FUNCTION fnc_peer_task_xp()
RETURNS TABLE
(
  Peer VARCHAR,
  Task VARCHAR,
  XP BIGINT
)
LANGUAGE SQL
AS $$
SELECT c.peer, SPLIT_PART(c.task, '_', 1) AS task, XP.xpamount
FROM Checks c
JOIN XP
  ON c.ID = XP."Check"
JOIN P2P
  ON c.ID = p2p."Check"
WHERE
  P2P.State = 'Success'::CHECKSTATUS
  AND (0 = (
    SELECT COUNT(*)
    FROM Verter
    WHERE "Check" = c.ID
      AND State = 'Start'::CHECKSTATUS)
  OR 1 = (
    SELECT COUNT(*)
    FROM Verter
    WHERE "Check" = c.ID
      AND State = 'Success'::CHECKSTATUS));
$$;

/* 3) Написать функцию, определяющую пиров, которые не выходили из кампуса в течение всего дня.
 Параметры функции: день, например 12.05.2022.
 Функция возвращает только список пиров.
*/
DROP FUNCTION IF EXISTS fnc_peer_stay_all_day_inside;
CREATE FUNCTION fnc_peer_stay_all_day_inside(p_date DATE)
RETURNS TABLE (Peer VARCHAR)
LANGUAGE SQL
AS $$
SELECT DISTINCT tt.Peer
FROM timetracking tt
WHERE
  tt."Date" = p_date
  AND (
        SELECT COUNT(*)
        FROM timetracking tt2
        WHERE
          tt2.Peer = tt.Peer
          AND tt2.State = '2'
          AND tt2."Date" = p_date
      ) = 1
$$;

/* 4) Найти процент успешных и неуспешных проверок за всё время.
 Формат вывода: процент успешных, процент неуспешных
*/
DROP PROCEDURE IF EXISTS prc_checks_success_rate;
CREATE PROCEDURE prc_checks_success_rate(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
OPEN "cursor" FOR
WITH cte AS (
  SELECT
    (
      SELECT COUNT(*)
      FROM P2P
      WHERE P2P.State = 'Success'::CHECKSTATUS
    ) AS p2psc,
    (
      SELECT COUNT(*)
      FROM Verter
      WHERE Verter.State = 'Failure'::CHECKSTATUS
    ) AS vf,
    (
      SELECT COUNT(*)
      FROM P2P
      WHERE P2P.State = 'Failure'::CHECKSTATUS
    ) AS p2pf
)
SELECT
  ROUND((p2psc - vf) / (p2psc + p2pf)::numeric * 100.0, 2)
  AS SuccessfulChecks,
  ROUND((p2pf + vf) / (p2psc + p2pf)::numeric * 100.0, 2)
  AS UnsuccessfulChecks
FROM cte;
END;
$$ LANGUAGE plpgsql;

/* 5) Посчитать изменение в количестве пир поинтов каждого пира по таблице TransferredPoints.
 Результат вывести отсортированным по изменению числа поинтов.
 Формат вывода:
    ник пира,
    изменение в количество пир поинтов
*/
DROP PROCEDURE IF EXISTS prc_total_change_prp_transferred;
CREATE PROCEDURE prc_total_change_prp_transferred(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
OPEN "cursor" FOR
WITH cte1 AS (
  SELECT
    checkingpeer AS p,
    SUM(pointsamount) AS s
  FROM transferredpoints
  GROUP BY checkingpeer
),
cte2 AS (
  SELECT
    checkedpeer AS p,
    SUM(pointsamount) AS s
  FROM transferredpoints
  GROUP BY checkedpeer
)
SELECT
  COALESCE(cte1.p, cte2.p) AS Peer,
  COALESCE(cte1.s, 0) - COALESCE(cte2.s, 0) AS PointsChange
FROM cte1
FULL JOIN cte2
  ON cte1.p = cte2.p
ORDER BY PointsChange;
END;
$$ LANGUAGE plpgsql;

/* 6) Посчитать изменение в количестве пир поинтов каждого пира по таблице,
 возвращаемой первой функцией из Part 3.
 Результат вывести отсортированным по изменению числа поинтов.
 Формат вывода:
    ник пира,
    изменение в количество пир поинтов
*/
DROP PROCEDURE IF EXISTS prc_total_change_prp_with_fnc;
CREATE PROCEDURE prc_total_change_prp_with_fnc(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
OPEN "cursor" FOR
WITH cte1 AS (
  SELECT
    peer1 AS p,
    SUM(PointsAmount) AS s
  FROM fnc_transferred_points_readable()
  GROUP BY peer1
),
cte2 AS (
  SELECT
    peer2 AS p,
    SUM(PointsAmount) AS s
  FROM fnc_transferred_points_readable()
  GROUP BY peer2
)
SELECT
  COALESCE(cte1.p, cte2.p) AS Peer,
  COALESCE(cte1.s, 0) - COALESCE(cte2.s, 0) AS PointsChange
FROM cte1
FULL JOIN cte2
  ON cte1.p = cte2.p
ORDER BY PointsChange;
END;
$$ LANGUAGE plpgsql;

/* 7) Определить самое часто проверяемое задание за каждый день.
 При одинаковом количестве проверок каких-то заданий в определенный день, вывести их все.
 Формат вывода:
    день,
    название задания
*/
DROP PROCEDURE IF EXISTS prc_most_checks_task_for_day;
CREATE PROCEDURE prc_most_checks_task_for_day(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
OPEN "cursor" FOR
WITH cte AS (
  SELECT
    Task,
    "Date",
    COUNT(*) AS cnt
  FROM Checks
  GROUP BY Task, "Date"
)
SELECT
  c."Date" AS Day,
  SPLIT_PART(c.Task, '_', 1) AS Task
FROM cte c
WHERE cnt = (
              SELECT MAX(c2.cnt)
              FROM cte c2
              WHERE c2."Date" = c."Date"
            );
END;
$$ LANGUAGE plpgsql;

/* 8) Определить длительность последней P2P проверки.
 Под длительностью подразумевается разница между временем, указанным в записи
 со статусом "начало", и временем, указанным в записи со статусом "успех" или "неуспех".
 Формат вывода: длительность проверки
*/
DROP PROCEDURE IF EXISTS prc_last_p2p_duration;
CREATE PROCEDURE prc_last_p2p_duration(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
OPEN "cursor" FOR
WITH cte AS (
    SELECT MAX("Check") AS max_id
    FROM p2p
    WHERE State IN ('Failure'::CHECKSTATUS, 'Success'::CHECKSTATUS)
)
SELECT
  (
    (
      SELECT "Time"
      FROM p2p
      WHERE State IN ('Failure'::CHECKSTATUS, 'Success'::CHECKSTATUS)
      AND "Check" = (SELECT max_id FROM cte)
    )
    -
    (
      SELECT "Time"
      FROM p2p
      WHERE State = 'Start'::CHECKSTATUS
      AND "Check" = (SELECT max_id FROM cte)
    )
  )::TIME AS CheckDuration;
END;
$$ LANGUAGE plpgsql;

/* 9) Найти всех пиров, выполнивших весь заданный блок задач и дату завершения последнего задания.
 Результат вывести отсортированным по дате завершения. 
 Формат вывода: 
    ник пира, 
    дата завершения блока (т.е. последнего выполненного задания из этого блока).
*/
DROP PROCEDURE IF EXISTS find_peers_completed_task_block;
CREATE PROCEDURE find_peers_completed_task_block(block_name VARCHAR, 
                                                 INOUT "cursor" REFCURSOR) 
AS $$ 
BEGIN 
  OPEN "cursor" FOR
  SELECT 
    checks.peer,
    MAX(checks."Date") AS "Day"
  FROM checks
  WHERE checks.peer 
    IN 
      (
        SELECT checks.peer 
        FROM checks 
        INNER JOIN p2p 
          ON checks.id = p2p."Check" 
          AND state = 'Success' 
          AND checks.task SIMILAR TO block_name || '\d+\_%'
        INNER JOIN verter 
          ON checks.id = verter."Check" 
          AND verter.state = 'Success' 
        GROUP BY checks.peer 
        HAVING COUNT(*) = (
                            SELECT COUNT(title) 
                            FROM tasks 
                            WHERE title SIMILAR TO block_name || '\d+\_%'
                          )
      )
  GROUP BY checks.peer
  ORDER BY "Day" DESC; 
END; 
$$ LANGUAGE plpgsql;

/* 10) Определить, к какому пиру стоит идти на проверку каждому обучающемуся.
 Определять нужно исходя из рекомендаций друзей пира, т.е. нужно найти пира, проверяться у которого рекомендует наибольшее число друзей. 
 Формат вывода: 
    ник пира, 
    ник найденного проверяющего.
*/
DROP PROCEDURE IF EXISTS find_recommended_peer_for_check;
CREATE PROCEDURE find_recommended_peer_for_check(INOUT "cursor" REFCURSOR) 
AS $$
BEGIN
  OPEN "cursor" FOR
  SELECT 
    peer1, 
    recommendedpeer
  FROM 
    (
      SELECT 
        peer1, 
        recommendedpeer, 
        count, 
        MAX(recommendations.count) OVER (PARTITION BY peer1) AS max_count
      FROM 
        (
          SELECT 
              peer1, 
              recommendedpeer, 
              COUNT(*) AS count 
          FROM friends 
          INNER JOIN recommendations 
            ON friends.peer2 = recommendations.peer
            GROUP BY peer1, recommendedpeer
        ) AS recommendations
    ) AS grouped_recommendations
  WHERE 
    grouped_recommendations.count = grouped_recommendations.max_count;
END;
$$ LANGUAGE plpgsql;

/* 11) Определить процент пиров, которые:
    а) Приступили только к блоку 1
    б) Приступили только к блоку 2
    с) Приступили к обоим
    д) Не приступили ни к одному
 Пир считается приступившим к блоку, если он проходил хоть одну проверку любого задания из этого блока (по таблице Checks).
 Формат вывода: 
    процент приступивших только к первому блоку, 
    процент приступивших только ко второму блоку, 
    процент приступивших к обоим, 
    процент не приступивших ни к одному
*/
DROP PROCEDURE IF EXISTS define_peers_persantage_in_tasks_blocks;
CREATE PROCEDURE 
  define_peers_persantage_in_tasks_blocks(block1 VARCHAR, 
                                          block2 VARCHAR, 
                                          INOUT "cursor" REFCURSOR)
AS $$ 
BEGIN
  OPEN "cursor" FOR
  WITH
  started_block1 AS (
    SELECT DISTINCT peer 
    FROM checks 
    WHERE task SIMILAR TO block1 || '\d+\_%'
  ),
  started_block2 AS (
    SELECT DISTINCT peer 
    FROM checks 
    WHERE task SIMILAR TO block2 || '\d+\_%'
  ),
  started_only_block1 AS (
    SELECT * FROM started_block1
    EXCEPT
    SELECT * FROM started_block2
  ),
  started_only_block2 AS (
    SELECT * FROM started_block2
    EXCEPT
    SELECT * FROM started_block1
  ),
  started_both_blocks AS (
    SELECT * FROM started_block1
    INTERSECT
    SELECT * FROM started_block2
  ),
  didnt_start_any_block AS (
    SELECT nickname FROM peers
    EXCEPT
    (
      SELECT * FROM started_block1
      UNION
      SELECT * FROM started_block2
    )
  )
  SELECT 
    (SELECT (COUNT(*) * 100.0 / total_peers.count)::BIGINT 
     FROM started_only_block1) AS "StartedBlock1",
    (SELECT (COUNT(*) * 100.0 / total_peers.count)::BIGINT 
     FROM started_only_block2) AS "StartedBlock2",
    (SELECT (COUNT(*) * 100.0 / total_peers.count)::BIGINT 
     FROM started_both_blocks) AS "StartedBothBlocks",
    (SELECT (COUNT(*) * 100.0 / total_peers.count)::BIGINT 
     FROM didnt_start_any_block) AS "DidntStartAnyBlock"
  FROM
    (SELECT COUNT(*) AS count FROM peers) AS total_peers;
END;
$$ LANGUAGE plpgsql;

/* 12) Определить N пиров с наибольшим числом друзей.
 Результат вывести отсортированным по кол-ву друзей. 
 Формат вывода:
    ник пира, 
    количество друзей.
*/
DROP PROCEDURE IF EXISTS find_N_peers_with_max_friends;
CREATE PROCEDURE find_N_peers_with_max_friends(peers_count BIGINT, 
                                               INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  SELECT 
    peers.nickname AS "Peer", 
    COUNT(friend) AS "FriendsCount"
  FROM 
    (
      (SELECT peer1 AS peer, peer2 AS friend FROM friends)
      UNION
      (SELECT peer2 AS peer, peer1 AS friend FROM friends)
    ) AS friends_union
  RIGHT JOIN peers 
    ON friends_union.peer = peers.nickname
  GROUP BY peers.nickname 
  ORDER BY "FriendsCount" DESC
  LIMIT peers_count;
END;
$$ LANGUAGE plpgsql;

/* 13) Определить процент пиров, которые когда-либо успешно проходили проверку в свой день рождения.
 Также определите процент пиров, которые хоть раз проваливали проверку в свой день рождения. 
 Формат вывода: 
    процент успехов в день рождения, 
    процент неуспехов в день рождения.
*/
DROP PROCEDURE IF EXISTS define_birthday_checks_statistics;
CREATE PROCEDURE define_birthday_checks_statistics(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  SELECT 
    (successful_checks * 100.0 / total_checks)::BIGINT AS "SuccessfulChecks", 
    ((total_checks - successful_checks) * 100.0 / total_checks)::BIGINT AS "UnsuccessfulChecks"
  FROM 
    (
      SELECT 
        COUNT(*) AS total_checks, 
        COUNT(xpamount) AS successful_checks 
      FROM checks
      INNER JOIN peers
        ON checks.peer = peers.nickname
        AND TO_CHAR(checks."Date", 'MM-DD') = TO_CHAR(peers.birthday, 'MM-DD')
      LEFT JOIN xp
        ON checks.id = xp."Check"
    ) AS birthday_checks;
END;
$$ LANGUAGE plpgsql;

/* 14) Определить кол-во XP, полученное в сумме каждым пиром.
 Если одна задача выполнена несколько раз, полученное за нее кол-во XP равно максимальному за эту задачу. 
 Результат вывести отсортированным по кол-ву XP. 
 Формат вывода: 
    ник пира, 
    количество XP.
*/
DROP PROCEDURE IF EXISTS define_total_received_xp_by_each_peer;
CREATE PROCEDURE define_total_received_xp_by_each_peer(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  SELECT peer, SUM(xp_per_block.max_xp) AS XP 
  FROM 
    (
      SELECT peer, task, MAX(xpamount) AS max_xp 
      FROM checks
      INNER JOIN xp 
        ON checks.id = xp."Check"
      GROUP BY peer, task) AS xp_per_block
  GROUP BY peer
  ORDER BY XP DESC;
END;
$$ LANGUAGE plpgsql;

/* 15) Определить всех пиров, которые сдали заданные задания 1 и 2, но не сдали задание 3.
 Формат вывода: список пиров.
*/
DROP PROCEDURE IF EXISTS define_peers_by_tasks;
CREATE PROCEDURE define_peers_by_tasks(compl_task1 VARCHAR,
                                       compl_task2 VARCHAR,
                                       uncoml_task VARCHAR,
                                       INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH completed_tasks AS (
    SELECT peer, task FROM checks
    INNER JOIN xp 
      ON checks.id = xp."Check"
  )
  (SELECT peer FROM completed_tasks WHERE task = compl_task1)
  INTERSECT
  (SELECT peer FROM completed_tasks WHERE task = compl_task2)
  EXCEPT
  (SELECT peer FROM completed_tasks WHERE task = uncoml_task)
  ORDER BY peer;
END;
$$ LANGUAGE plpgsql;

/* 16) Используя рекурсивное обобщенное табличное выражение, для каждой задачи вывести кол-во предшествующих ей задач.
 То есть сколько задач нужно выполнить, исходя из условий входа, чтобы получить доступ к текущей. 
 Формат вывода: 
    название задачи, 
    количество предшествующих.
*/
DROP PROCEDURE IF EXISTS define_parent_tasks_count;
CREATE PROCEDURE define_parent_tasks_count(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH RECURSIVE cte_tasks AS (
    SELECT
      title,
      parenttask
    FROM tasks
    UNION ALL
    SELECT
      cte_tasks.title,
      tasks.parenttask
    FROM cte_tasks 
    INNER JOIN tasks 
      ON cte_tasks.parenttask = tasks.title
  )
  SELECT 
    title AS Task, 
    COUNT (DISTINCT parenttask) AS PrevCount 
  FROM cte_tasks
  GROUP BY title;
END;
$$ LANGUAGE plpgsql;

/* 17) Найти "удачные" для проверок дни. 
 День считается "удачным", если в нем есть хотя бы N идущих подряд успешных проверки.
 Временем проверки считать время начала P2P этапа. 
 Под идущими подряд успешными проверками подразумеваются успешные проверки, между которыми нет неуспешных. 
 При этом кол-во опыта за каждую из этих проверок должно быть не меньше 80% от максимального. 
 Формат вывода: список дней.
*/
DROP PROCEDURE IF EXISTS find_lucky_days_for_checks;
CREATE PROCEDURE find_lucky_days_for_checks(
	consecutive_checks BIGINT, INOUT "cursor" REFCURSOR
)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH rawdata AS (
    SELECT
      Checks."Date" AS checkdate,
      P2P."Time" AS checktime,
      XP.xpamount,
      Tasks.maxxp
    FROM Checks
    INNER JOIN P2P
      ON Checks.id = P2P."Check" 
      AND P2P.state = 'Start'
    LEFT JOIN XP
      ON Checks.id = XP."Check"
    INNER JOIN Tasks
      ON Checks.task = Tasks.title
    ORDER BY checkdate, checktime
  ),
  addrnum AS (
    SELECT
      row_number() OVER () AS rnum,
      *
    FROM rawdata
  ),
  partitioned AS (
    SELECT
      *,
      rnum - ROW_NUMBER() OVER (PARTITION BY checkdate) AS grp
    FROM addrnum
    WHERE xpamount >= (maxxp * 0.8)
  ),
  counted AS (
    SELECT
    *,
    COUNT (*) OVER (PARTITION BY checkdate, grp) AS cnt
    FROM partitioned
  )
  SELECT DISTINCT checkdate AS luckyday 
  FROM counted 
  WHERE cnt >= consecutive_checks;
END;
$$ LANGUAGE plpgsql;

/* 18) Определить пира с наибольшим числом выполненных заданий
 Формат вывода: 
    ник пира, 
    число выполненных заданий
*/
DROP PROCEDURE IF EXISTS find_peer_with_max_completed_tasks;
CREATE PROCEDURE find_peer_with_max_completed_tasks(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH completed_tasks AS (
    SELECT
      Checks.peer,
      Checks.task
    FROM Checks
    INNER JOIN XP
      ON Checks.id = XP."Check"
    GROUP BY Checks.peer, Checks.task
  ),
  counted_tasks AS (
    SELECT
      peer,
      COUNT (*) AS tasks_count
    FROM completed_tasks
    GROUP BY peer
  )
  SELECT * FROM counted_tasks
  ORDER BY tasks_count DESC, peer
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

/* 19) Определить пира с наибольшим количеством XP
 Формат вывода: 
    ник пира, 
    количество XP
*/
DROP PROCEDURE IF EXISTS find_peer_with_max_xp;
CREATE PROCEDURE find_peer_with_max_xp(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH xp_per_tasks AS (
    SELECT peer, task, MAX(xpamount) AS xpgain 
	FROM Checks
	INNER JOIN XP 
	ON Checks.id = XP."Check"
	GROUP BY peer, task
  ),
  total_exp AS (
    SELECT 
      peer, 
      SUM(xpgain) AS xp 
    FROM xp_per_tasks
    GROUP BY peer
  )
  SELECT * FROM total_exp
  ORDER BY xp DESC, peer
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

/* 20) Определить пира, который провел сегодня в кампусе больше всего времени
 Формат вывода: ник пира
*/
DROP PROCEDURE IF EXISTS find_peer_with_max_staytime_for_today;
CREATE PROCEDURE find_peer_with_max_staytime_for_today(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH raw AS (
    SELECT
      peer,
      "Date" AS visitdate,
      "Time" AS visittime,
      state
    FROM TimeTracking 
    ORDER BY Peer, "Date", "Time", state
  ),
  arrivals AS (
    SELECT
      row_number() OVER () AS rnum,
      *
    FROM raw WHERE state = 1
  ),
  departures AS (
    SELECT
      row_number() OVER () AS rnum,
      *
    FROM raw WHERE state = 2
  ),
  visits AS (
    SELECT
      arrivals.peer,
      arrivals.visitdate,
      arrivals.visittime AS arrivaltime,
      departures.visittime AS departuretime
    FROM arrivals
    INNER JOIN departures
      ON arrivals.rnum = departures.rnum
  ),
  staytime AS (
    SELECT
    peer,
    visitdate,
    SUM(departuretime - arrivaltime) AS staytime
    FROM visits
    GROUP BY peer, visitdate
  )
  SELECT peer
  FROM staytime
  WHERE visitdate = CURRENT_DATE
  ORDER BY staytime DESC
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

/* 21) Определить пиров, приходивших раньше заданного времени не менее N раз за всё время
 Формат вывода: список пиров
*/
DROP PROCEDURE IF EXISTS find_peers_with_N_early_visits;
CREATE PROCEDURE find_peers_with_N_early_visits(
  earlytime TIME, visitnum BIGINT, INOUT "cursor" REFCURSOR
)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH
  arrivals AS (
    SELECT
      peer,
      "Date" AS visitdate,
      "Time" AS visittime
    FROM TimeTracking 
    WHERE state = 1
  ),
  earliest_arrivals AS (
    SELECT
      peer,
      visitdate,
      MIN(visittime) AS firstvisit
    FROM arrivals
    GROUP BY peer, visitdate
  ),
  visit_count AS (
    SELECT
      peer,
      COUNT (*) AS cnt
    FROM earliest_arrivals
    WHERE firstvisit < earlytime
    GROUP BY peer
  )
  SELECT peer
  FROM visit_count
  WHERE cnt >= visitnum;
END;
$$ LANGUAGE plpgsql;

/* 22) Определить пиров, выходивших за последние N дней из кампуса больше M раз
 Формат вывода: список пиров
*/
DROP PROCEDURE IF EXISTS find_peers_who_left_campus_more_than_M_times_during_N_days;
CREATE PROCEDURE 
  find_peers_who_left_campus_more_than_M_times_during_N_days(deptnum BIGINT, 
                                                             daynum BIGINT, 
                                                             INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH
  departures AS (
    SELECT *
    FROM TimeTracking
    WHERE state = 2 
      AND "Date" BETWEEN (CURRENT_DATE - daynum + 1) AND CURRENT_DATE
  ),
  counted AS (
    SELECT
      peer,
      COUNT (*) AS cnt
    FROM departures
    GROUP BY peer
  )
  SELECT peer 
  FROM counted
  WHERE cnt > deptnum;
END;
$$ LANGUAGE plpgsql;

/* 23) Определить пира, который пришел сегодня последним
 Формат вывода: список пиров
*/
DROP PROCEDURE IF EXISTS find_peer_which_arrived_last;
CREATE PROCEDURE find_peer_which_arrived_last(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH
  arrivals AS (
    SELECT
      peer,
      "Date" AS visitdate,
      "Time" AS visittime
    FROM TimeTracking 
    WHERE state = 1 AND "Date" = CURRENT_DATE
  ),
  earliest_arrivals AS (
    SELECT
      peer,
      visitdate,
      MIN(visittime) AS firstvisit
    FROM arrivals
    GROUP BY peer, visitdate
  )
  SELECT peer
  FROM earliest_arrivals
  ORDER BY firstvisit DESC
  LIMIT 1;
END;
$$ LANGUAGE plpgsql;

/* 24) Определить пиров, которые выходили вчера из кампуса больше чем на N минут
 Формат вывода: список пиров
*/
DROP PROCEDURE IF EXISTS find_peers_who_left_campus_for_N_minutes_yesterday;
CREATE PROCEDURE 
  find_peers_who_left_campus_for_N_minutes_yesterday(minnum BIGINT, 
                                                     INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH
  raw AS (
    SELECT
      peer,
      "Date" AS visitdate,
      "Time" AS visittime,
      state
    FROM TimeTracking
    WHERE "Date" = (CURRENT_DATE - 1)
    ORDER BY Peer, "Date", "Time", state
  ),
  arrivals AS (
    SELECT
      row_number() OVER () AS rnum,
      *
    FROM raw WHERE state = 1
  ),
  departures AS (
    SELECT
      row_number() OVER () AS rnum,
      *
    FROM raw WHERE state = 2
  ),
  earliest_arrivals AS (
    SELECT
      peer,
      visitdate,
      MIN(visittime) AS firstvisit
    FROM arrivals
    GROUP BY peer, visitdate
  ),
  latest_departures AS (
    SELECT
      peer,
      visitdate,
      MAX(visittime) AS lastdepart
    FROM departures
    GROUP BY peer, visitdate
  ),
  visits AS (
    SELECT
      arrivals.peer,
      arrivals.visitdate,
      arrivals.visittime AS arrivaltime,
      departures.visittime AS departuretime
    FROM arrivals
    INNER JOIN departures
      ON arrivals.rnum = departures.rnum
  ),
  earliest_to_latest AS (
    SELECT
      earliest_arrivals.peer,
      earliest_arrivals.visitdate,
      firstvisit,
      lastdepart
    FROM earliest_arrivals
    INNER JOIN latest_departures
      ON earliest_arrivals.peer = latest_departures.peer
      AND earliest_arrivals.visitdate = latest_departures.visitdate
  ),
  staytime AS (
    SELECT
      peer,
      visitdate,
      SUM(departuretime - arrivaltime) AS staytime
    FROM visits
    GROUP BY peer, visitdate
  ),
  absenttime AS (
    SELECT
      staytime.peer,
      staytime.visitdate,
      staytime,
      firstvisit,
      lastdepart,
      (lastdepart - firstvisit - staytime) AS absenttime
    FROM staytime
    INNER JOIN earliest_to_latest
      ON staytime.peer = earliest_to_latest.peer
      AND staytime.visitdate = earliest_to_latest.visitdate
  )
  SELECT peer
  FROM absenttime
  WHERE absenttime > MAKE_INTERVAL(mins := minnum);
END;
$$ LANGUAGE plpgsql;

/* 25) Определить для каждого месяца процент ранних входов
 Для каждого месяца посчитать, сколько раз люди, родившиеся в этот месяц, 
 приходили в кампус за всё время (будем называть это общим числом входов). 
 Для каждого месяца посчитать, сколько раз люди, родившиеся в этот месяц, 
 приходили в кампус раньше 12:00 за всё время (будем называть это числом ранних входов). 
 Для каждого месяца посчитать процент ранних входов в кампус относительно общего числа входов. 
 Формат вывода: 
    месяц, 
    процент ранних входов
*/
DROP PROCEDURE IF EXISTS find_early_entries_percentage_per_month;
CREATE PROCEDURE find_early_entries_percentage_per_month(INOUT "cursor" REFCURSOR)
AS $$
BEGIN
  OPEN "cursor" FOR
  WITH
  monthrange AS (
    SELECT 
      month_id,
      TO_CHAR(TO_DATE(month_id::text, 'MM'), 'Month') AS monthname
    FROM GENERATE_SERIES(1, 12) AS g(month_id)
  ),
  peers_per_month AS (
    SELECT
      *
    FROM monthrange
    LEFT JOIN Peers
      ON monthrange.month_id = DATE_PART('month', Peers.birthday)
    ORDER BY month_id
  ),
  arrivals AS (
    SELECT
      peer,
      "Date" AS visitdate,
      "Time" AS visittime
    FROM TimeTracking 
    WHERE state = 1
  ),
  earliest_arrivals AS (
    SELECT
      peer,
      visitdate,
      MIN(visittime) AS firstvisit
    FROM arrivals
    GROUP BY peer, visitdate
  ),
  grouped AS (
    SELECT
      peer,
      COUNT (*) FILTER (WHERE firstvisit < '12:00:00') AS earlyvisitcnt,
      COUNT (*) AS visitcnt
    FROM earliest_arrivals
    GROUP BY peer
  ),
  aggregated AS (
    SELECT
      month_id,
      monthname,
      SUM(earlyvisitcnt) AS earlytotal,
      SUM(visitcnt) AS total
    FROM peers_per_month
    LEFT JOIN grouped
      ON peers_per_month.nickname = grouped.peer
    GROUP BY month_id, monthname
    ORDER BY month_id
  )
  SELECT
    monthname AS month,
    COALESCE(ROUND(earlytotal / total * 100), 0) AS EarlyEntries
  FROM aggregated;
END;
$$ LANGUAGE plpgsql;


