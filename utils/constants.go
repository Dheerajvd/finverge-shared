package utils

const SuperAdmin = "67a62567da7b3ebfbfcea4b01"
const Admin = "67a6256c7a73ebfbfcea4b12"
const CreditOps = "67a62563da7b3ebfbfcea4b23"
const DisburseOps = "67a62567da7b3ebfbfcea4b34"
const AuditOps = "67a62567da7b3ebfbfcea4b45"
const Customer = "67a62567da7b3ebfbfcea4b56"

const DefaultRole = "default"

func GetroleId(role string) string {
	switch role {
	case "super-admin":
		return SuperAdmin
	case "admin":
		return Admin
	case "credit-ops":
		return CreditOps
	case "disburse-ops":
		return DisburseOps
	case "audit-ops":
		return AuditOps
	case "customer":
		return DefaultRole
	default:
		return DefaultRole
	}
}

func GetRoleName(roleId string) string {
	switch roleId {
	case SuperAdmin:
		return "super-admin"
	case Admin:
		return "admin"
	case CreditOps:
		return "credit-ops"
	case DisburseOps:
		return "disbures-ops"
	case AuditOps:
		return "audit-ops"
	case Customer:
		return "customer"
	default:
		return DefaultRole
	}
}
