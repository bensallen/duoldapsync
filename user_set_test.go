package main

import (
	"net/url"
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

func TestUser_URLValues(t *testing.T) {
	type fields struct {
		FullName    string
		Email       string
		FirstName   string
		LastName    string
		LDAP        bool
		Duo         bool
		NeedsUpdate bool
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "Username only",
			args: args{username: "test1"},
			want: url.Values{"username": []string{"test1"}},
		},
		{
			name:    "No username",
			args:    args{},
			wantErr: true,
		},
		{
			name:   "Username only",
			args:   args{username: "test1"},
			fields: fields{FullName: "Test One", Email: "test@example.com", FirstName: "Test", LastName: "One"},
			want: url.Values{
				"username":  []string{"test1"},
				"realname":  []string{"Test One"},
				"email":     []string{"test@example.com"},
				"firstname": []string{"Test"},
				"lastname":  []string{"One"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				FullName:    tt.fields.FullName,
				Email:       tt.fields.Email,
				FirstName:   tt.fields.FirstName,
				LastName:    tt.fields.LastName,
				LDAP:        tt.fields.LDAP,
				Duo:         tt.fields.Duo,
				NeedsUpdate: tt.fields.NeedsUpdate,
			}
			got, err := u.URLValues(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.URLValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
