package queries

const GetFunctionData = `SELECT
	pg_proc.oid AS id,
	proname::text AS name,
	(SELECT lanname FROM pg_language WHERE oid = prolang)::text AS lngname,
	prokind::text AS kind,
	pronargs AS argnumber,
	(SELECT typname FROM pg_type WHERE oid = prorettype)::text AS returntype,
	ARRAY(SELECT (SELECT typname FROM pg_type WHERE oid = t) FROM unnest(proargtypes) AS t)::text[] AS inargs,
	ARRAY(SELECT (SELECT typname FROM pg_type WHERE oid = t) FROM unnest(proallargtypes) AS t)::text[] AS allargs,
	proargmodes::text[] AS argmodes,
	proargnames AS argnames
FROM pg_proc
JOIN pg_namespace ON pg_proc.pronamespace = pg_namespace.oid
WHERE nspname = 'public'
`
