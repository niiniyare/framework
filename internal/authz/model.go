package authz

// CasbinModel is the RBAC model definition for the Awo framework.
// Uses three-parameter domain-scoped roles with deny-override policy effect.
const CasbinModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = (g(r.sub, p.sub, r.dom) || p.sub == "*") && r.dom == p.dom && keyMatch2(r.obj, p.obj) && keyMatch(r.act, p.act)
`
