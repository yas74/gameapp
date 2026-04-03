package mysqlaccesscontrol

import (
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/errmsg"
	"gocasts/gameapp/pkg/richerror"
	"gocasts/gameapp/pkg/slice"
	"gocasts/gameapp/repository/mysql"
	"strings"
	"time"
)

func (d *DB) GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error) {
	const op = "accesscontrol.GetUserPermissionTitles"

	// user, err := d.GetUserByID(userID)
	// if err != nil {
	// 	return nil, richerror.New(op).WithErr(err)
	// }

	roleAcl := make([]entity.AccessControl, 0)
	rows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ? and actor_id = ?`, entity.RoleActorType, role)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}

		roleAcl = append(roleAcl, acl)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	userAcl := make([]entity.AccessControl, 0)

	userRows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ? and actor_id = ?`, entity.UserActorType, userID)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}

		userAcl = append(userAcl, acl)
	}

	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	// merge acl by permission id
	permissionIDs := make([]uint, 0)
	for _, r := range roleAcl {
		if !slice.DoesExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	args := make([]interface{}, len(permissionIDs))

	for i, id := range permissionIDs {
		args[i] = id
	}

	query := "select * from permissions where id in (?" +
		strings.Repeat(",?", len(permissionIDs)-1) +
		")"

	pRows, err := d.conn.Conn().Query(
		query,
		args...,
	)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	defer pRows.Close()

	permissionTitles := make([]entity.PermissionTitle, 0)

	for pRows.Next() {
		permission, err := scanPermission(pRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
		}
		permissionTitles = append(permissionTitles, permission.Title)
	}

	if err := pRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return permissionTitles, nil
	// for _, r := range userAcl {
	// 	if !slice.DoesExist(permissionIDs, r.PermissionID) {
	// 		permissionIDs = append(permissionIDs, r.PermissionID)
	// 	}
	// }

	// for i := range permissionIDs {
	// 	permission :=
	// }
}

func scanAccessControl(scanner mysql.Scanner) (entity.AccessControl, error) {
	var createdAt time.Time
	acl := entity.AccessControl{}

	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)

	return acl, err
}
