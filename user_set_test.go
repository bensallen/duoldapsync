package main

import (
	"reflect"
	"testing"

	ldap "gopkg.in/ldap.v2"
)

func TestUserSet_AddLDAPEntries(t *testing.T) {
	type args struct {
		entries        []*ldap.Entry
		ldapUserSearch *LDAPUserSearch
	}
	tests := []struct {
		name string
		u    UserSet
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.AddLDAPEntries(tt.args.entries, tt.args.ldapUserSearch)
		})
	}
}

func TestUserSet_AddDuoResults(t *testing.T) {
	type args struct {
		result *UsersResponse
	}

	tests := []struct {
		name  string
		u     UserSet
		args  args
		wants UserSet
	}{
		{
			name: "example1 Single User",
			args: args{
				result: &UsersResponse{
					Response: []UserResponse{
						UserResponse{
							Username: "example1"},
					},
				},
			},
			u:     UserSet{},
			wants: UserSet{"example1": &User{Duo: true}},
		},
		{
			name: "Existing User and New User",
			args: args{
				result: &UsersResponse{
					Response: []UserResponse{
						UserResponse{
							Username: "example2"},
					},
				},
			},
			u:     UserSet{"example1": &User{Duo: false, LDAP: true}},
			wants: UserSet{"example1": &User{Duo: false, LDAP: true}, "example2": &User{Duo: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.AddDuoResults(tt.args.result)
			if !reflect.DeepEqual(tt.u, tt.wants) {
				t.Fatalf("Mismatch between result %v and wants %v", tt.u, tt.wants)
			}
		})
	}
}
