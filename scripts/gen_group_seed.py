#!/usr/bin/env python3
"""
Transform etc/casdoor/init_data.json to use UUIDs as group names,
and update the pg_groups INSERT block in scripts/data/postgres/init.sql.

Casdoor group name = deterministic UUID (uuid5)
PostgreSQL pg_groups table stores UUID <-> display name mapping

Usage:
    python3 scripts/gen_group_seed.py
"""
import json
import re
import uuid
import os

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC      = os.path.join(BASE_DIR, 'etc', 'casdoor', 'init_data.json')
DST      = SRC  # overwrite in place
INIT_SQL = os.path.join(BASE_DIR, 'scripts', 'data', 'postgres', 'init.sql')

# Deterministic namespace for cmn-core group UUIDs
NAMESPACE = uuid.UUID('a1b2c3d4-e5f6-7890-abcd-ef1234567890')

_UUID_RE = re.compile(
    r'^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$', re.I)

_INSERT_RE = re.compile(
    r'INSERT INTO pg_groups \(uuid, name, created_at, updated_at\) VALUES\n'
    r'.*?'
    r'ON CONFLICT \(uuid\) DO UPDATE\n  SET name = EXCLUDED\.name, updated_at = NOW\(\);',
    re.DOTALL,
)


def group_uuid(display_name):
    return str(uuid.uuid5(NAMESPACE, display_name))


def display_name_from_group(g):
    if _UUID_RE.match(g.get('name', '')):
        dn = g.get('displayName', '')
        m = re.match(r'^Group\s+(\d+)$', dn, re.IGNORECASE)
        if m:
            return 'group{:03d}'.format(int(m.group(1)))
        return dn.lower().replace(' ', '')
    return g['name']


def main():
    with open(SRC) as f:
        data = json.load(f)

    groups_in_json = data.get('groups', [])

    if not groups_in_json:
        names = set()
        for u in data.get('users', []):
            for g in u.get('groups', []):
                names.add(g.split('/', 1)[-1])
        display_names = sorted(names)
    else:
        display_names = [display_name_from_group(g) for g in groups_in_json]

    name_to_uuid = {n: group_uuid(n) for n in display_names}
    print('  Groups to migrate: {}'.format(len(name_to_uuid)))

    for g in groups_in_json:
        if not _UUID_RE.match(g['name']):
            g['name'] = name_to_uuid[g['name']]

    for u in data.get('users', []):
        new_groups = []
        for g in u.get('groups', []):
            parts   = g.split('/', 1)
            display = parts[-1]
            prefix  = parts[0] + '/' if len(parts) == 2 else ''
            if _UUID_RE.match(display):
                new_groups.append(g)
            else:
                new_groups.append(prefix + name_to_uuid.get(display, display))
        u['groups'] = new_groups

        old_props = u.get('properties') or {}
        new_props = {}
        for k, v in old_props.items():
            replaced = False
            for display, uid in name_to_uuid.items():
                if k == 'cmn_group_{}_role'.format(display):
                    new_props['cmn_group_{}_role'.format(uid)] = v
                    replaced = True
                    break
            if not replaced:
                new_props[k] = v
        u['properties'] = new_props

    with open(DST, 'w') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
        f.write('\n')
    print('  Written: {}'.format(DST))

    # Build INSERT block
    rows_data = sorted(
        [(display_name_from_group(g), g['name']) for g in data.get('groups', [])],
        key=lambda x: x[0],
    )
    insert_rows = ',\n'.join(
        "  ('{}', '{}', NOW(), NOW())".format(uid, short)
        for short, uid in rows_data
    )
    insert_block = (
        'INSERT INTO pg_groups (uuid, name, created_at, updated_at) VALUES\n'
        + insert_rows + '\n'
        + 'ON CONFLICT (uuid) DO UPDATE\n'
        + '  SET name = EXCLUDED.name, updated_at = NOW();'
    )

    # Patch init.sql using string find+replace (avoids re.sub backslash issues)
    with open(INIT_SQL) as f:
        sql = f.read()

    m = _INSERT_RE.search(sql)
    if m is None:
        print('  WARNING: pg_groups INSERT block not found in init.sql')
    else:
        patched = sql[:m.start()] + insert_block + sql[m.end():]
        with open(INIT_SQL, 'w') as f:
            f.write(patched)
        print('  Updated: {}'.format(INIT_SQL))


if __name__ == '__main__':
    main()
