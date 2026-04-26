#!/usr/bin/env python3
"""Generate seed data: 50 users, 100 groups, 10 members/group for Casdoor and Keycloak."""
import json

NUM_USERS   = 50
NUM_GROUPS  = 100
MEMBERS_PER = 10

def user_name(i):  return f"user{i:02d}"
def group_name(g): return f"group{g+1:03d}"

def members_of(g):
    return [(g * MEMBERS_PER + k) % NUM_USERS + 1 for k in range(MEMBERS_PER)]

user_groups = {i: [] for i in range(1, NUM_USERS + 1)}
for g in range(NUM_GROUPS):
    for i in members_of(g):
        user_groups[i].append(g)

# ── Casdoor ──────────────────────────────────────────────────────────────────
with open("etc/casdoor/init_data.json") as f:
    casdoor = json.load(f)

casdoor["users"] = [
    {
        "owner": "cmn", "name": user_name(i), "type": "normal-user",
        "password": "Password123!", "displayName": f"User {i:02d}",
        "firstName": "User", "lastName": f"{i:02d}",
        "email": f"{user_name(i)}@cmn.local", "emailVerified": False,
        "isAdmin": False, "isGlobalAdmin": False,
        "signupApplication": "app-cmn-core", "score": 2000,
        "groups": [f"cmn/{group_name(g)}" for g in user_groups[i]],
    }
    for i in range(1, NUM_USERS + 1)
]

casdoor["groups"] = [
    {
        "owner": "cmn", "name": group_name(g), "displayName": f"Group {g+1:03d}",
        "type": "Virtual", "parentId": "", "isPublic": True,
        "members": [f"cmn/{user_name(i)}" for i in members_of(g)],
    }
    for g in range(NUM_GROUPS)
]

# Role for all app users → used to grant /api/get-groups permission
casdoor["roles"] = [
    {
        "owner": "cmn", "name": "role-app-user", "displayName": "App User",
        "description": "Regular application user",
        "users": [f"cmn/{user_name(i)}" for i in range(1, NUM_USERS + 1)],
        "groups": [], "roles": [], "domains": [], "isEnabled": True,
    }
]

# Without this, Casdoor profile page Groups field shows "No data":
# the TreeSelect calls GET /api/get-groups which is admin-only by default.
casdoor["permissions"] = [
    {
        "owner": "cmn", "name": "perm-app-user-read-groups",
        "displayName": "App users can list groups",
        "description": "Allows GET /api/get-groups so profile Groups field renders",
        # List users directly (belt-and-suspenders: works even if role expansion fails)
        "users": [f"cmn/{user_name(i)}" for i in range(1, NUM_USERS + 1)],
        "groups": [], "roles": ["cmn/role-app-user"],
        "domains": [], "model": "", "adapter": "",
        "resourceType": "URL", "resources": ["/api/get-groups"],
        "actions": ["GET"], "effect": "Allow", "isEnabled": True,
    }
]

with open("etc/casdoor/init_data.json", "w") as f:
    json.dump(casdoor, f, indent=2, ensure_ascii=False)
    f.write("\n")

# ── Keycloak ──────────────────────────────────────────────────────────────────
with open("etc/keycloak/cmn-realm.json") as f:
    keycloak = json.load(f)

keycloak["users"] = [
    {
        "username": user_name(i), "email": f"{user_name(i)}@cmn.local",
        "enabled": True, "firstName": "User", "lastName": f"{i:02d}",
        "credentials": [{"type": "password", "value": "Password123!", "temporary": False}],
        "realmRoles": ["app"],
        "groups": [f"/{group_name(g)}" for g in user_groups[i]],
    }
    for i in range(1, NUM_USERS + 1)
]

keycloak["groups"] = [
    {
        "name": group_name(g), "path": f"/{group_name(g)}",
        "attributes": {}, "realmRoles": [], "clientRoles": {}, "subGroups": [],
    }
    for g in range(NUM_GROUPS)
]

with open("etc/keycloak/cmn-realm.json", "w") as f:
    json.dump(keycloak, f, indent=2, ensure_ascii=False)
    f.write("\n")

# ── verify ────────────────────────────────────────────────────────────────────
print(f"casdoor  : {len(casdoor['users'])} users, {len(casdoor['groups'])} groups")
print(f"           {len(casdoor['roles'])} roles, {len(casdoor['permissions'])} permissions")
print(f"keycloak : {len(keycloak['users'])} users, {len(keycloak['groups'])} groups")
perm = casdoor["permissions"][0]
print(f"permission '{perm['name']}': {perm['resources']} → {perm['effect']}")
