// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	pb "github.com/casbin/casbin-server/proto"
)

// GetRolesForUser gets the roles that a user has.
func (s *Server) GetRolesForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	res, _ := e.GetModel()["g"]["g"].RM.GetRoles(in.User)

	return &pb.ArrayReply{Array: res}, nil
}

// GetUsersForRole gets the users that has a role.
func (s *Server) GetUsersForRole(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	res, _ := e.GetModel()["g"]["g"].RM.GetUsers(in.User)

	return &pb.ArrayReply{Array: res}, nil
}

// HasRoleForUser determines whether a user has a role.
func (s *Server) HasRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	roles := e.GetRolesForUser(in.User)

	for _, r := range roles {
		if r == in.Role {
			return &pb.BoolReply{Res: true}, nil
		}
	}

	return &pb.BoolReply{}, nil
}

// AddRoleForUser adds a role for a user.
// Returns false if the user already has the role (aka not affected).
func (s *Server) AddRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.AddGroupingPolicy(in.User, in.Role)}, nil
}

// DeleteRoleForUser deletes a role for a user.
// Returns false if the user does not have the role (aka not affected).
func (s *Server) DeleteRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemoveGroupingPolicy(in.User, in.Role)}, nil
}

// DeleteRolesForUser deletes all roles for a user.
// Returns false if the user does not have any roles (aka not affected).
func (s *Server) DeleteRolesForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemoveFilteredGroupingPolicy(0, in.User)}, nil
}

// DeleteUser deletes a user.
// Returns false if the user does not exist (aka not affected).
func (s *Server) DeleteUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemoveFilteredGroupingPolicy(0, in.User)}, nil
}

// DeleteRole deletes a role.
func (s *Server) DeleteRole(ctx context.Context, in *pb.UserRoleRequest) (*pb.EmptyReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	e.RemoveFilteredGroupingPolicy(1, in.Role)
	e.RemoveFilteredPolicy(0, in.Role)

	return &pb.EmptyReply{}, nil
}

// DeletePermission deletes a permission.
// Returns false if the permission does not exist (aka not affected).
func (s *Server) DeletePermission(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemoveFilteredPolicy(1, in.Permissions...)}, nil
}

// AddPermissionForUser adds a permission for a user or role.
// Returns false if the user or role already has the permission (aka not affected).
func (s *Server) AddPermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.AddPolicy(s.convertPermissions(in.User, in.Permissions...)...)}, nil
}

// DeletePermissionForUser deletes a permission for a user or role.
// Returns false if the user or role does not have the permission (aka not affected).
func (s *Server) DeletePermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemovePolicy(s.convertPermissions(in.User, in.Permissions...)...)}, nil
}

// DeletePermissionsForUser deletes permissions for a user or role.
// Returns false if the user or role does not have any permissions (aka not affected).
func (s *Server) DeletePermissionsForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.RemoveFilteredPolicy(0, in.User)}, nil
}

// GetPermissionsForUser gets permissions for a user or role.
func (s *Server) GetPermissionsForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetFilteredPolicy(0, in.User)), nil
}

// HasPermissionForUser determines whether a user has a permission.
func (s *Server) HasPermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.HasPolicy(s.convertPermissions(in.User, in.Permissions...)...)}, nil
}

func (s *Server) convertPermissions(user string, permissions ...string) []interface{} {
	params := make([]interface{}, 0, len(permissions)+1)
	params = append(params, user)
	for _, perm := range permissions {
		params = append(params, perm)
	}

	return params
}

// GetImplicitRolesForUser gets implicit roles that a user has.
// Compared to GetRolesForUser(), this function retrieves indirect roles besides direct roles.
// For example:
// g, alice, role:admin
// g, role:admin, role:user
//
// GetRolesForUser("alice") can only get: ["role:admin"].
// But GetImplicitRolesForUser("alice") will get: ["role:admin", "role:user"].
func (s *Server) GetImplicitRolesForUser(ctx context.Context, in *pb.ImplicitRolesForUserRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	return &pb.ArrayReply{Array: e.GetImplicitRolesForUser(in.User, in.Domain...)}, nil
}

// GetImplicitPermissionsForUser gets implicit permissions for a user or role.
// Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles.
// For example:
// p, admin, data1, read
// p, alice, data2, read
// g, alice, admin
//
// GetPermissionsForUser("alice") can only get: [["alice", "data2", "read"]].
// But GetImplicitPermissionsForUser("alice") will get: [["admin", "data1", "read"], ["alice", "data2", "read"]].
func (s *Server) GetImplicitPermissionsForUser(ctx context.Context, in *pb.ImplicitPermissionsForUserRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetImplicitPermissionsForUser(in.User)), nil
}

// HasImplicitPermissionForUser check implicit permission for a user
// Compared to GetPermissionsForUser(), this function retrieves permissions for inherited roles.
// For example:
// p, admin, data1, read
// p, alice, data2, read
// g, alice, admin
//
// HasImplicitPermissionForUser("alice", ["data1", "read"]) is true
// HasImplicitPermissionForUser("alice", ["data2", "read"]) is true
// HasImplicitPermissionForUser("admin", ["data1", "read"]) is true
// HasImplicitPermissionForUser("admin", ["data2", "read"]) is false
func (s *Server) HasImplicitPermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	matched := false
	for _, fieldValue := range e.GetImplicitPermissionsForUser(in.User) {
		if arrayEquals(fieldValue[1:], in.Permissions) {
			matched = true
			break
		}
	}

	return &pb.BoolReply{Res: matched}, nil
}

func arrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
