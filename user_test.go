package main

import (
	"net/url"
	"reflect"
	"testing"

	ldap "gopkg.in/ldap.v2"
)

func TestUserSet_addLDAPEntries(t *testing.T) {
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
			tt.u.addLDAPEntries(tt.args.entries, tt.args.ldapUserSearch)
		})
	}
}

func TestUserSet_addDuoResults(t *testing.T) {
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
						{
							Username: "example1"},
					},
				},
			},
			u:     UserSet{},
			wants: UserSet{"example1": &User{Duo: true, Username: "example1"}},
		},
		{
			name: "Existing User and New User",
			args: args{
				result: &UsersResponse{
					Response: []UserResponse{
						{
							Username: "example2"},
					},
				},
			},
			u:     UserSet{"example1": &User{Duo: false, LDAP: true, Username: "example1"}},
			wants: UserSet{"example1": &User{Duo: false, LDAP: true, Username: "example1"}, "example2": &User{Duo: true, Username: "example2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.addDuoResults(tt.args.result)
			if !reflect.DeepEqual(tt.u, tt.wants) {
				t.Fatalf("Mismatch between result %v and wants %v", tt.u, tt.wants)
			}
		})
	}
}

func TestUser_urlValues(t *testing.T) {
	type fields struct {
		Username    string
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
		want    url.Values
		wantErr bool
	}{
		{
			name:   "Username only",
			fields: fields{Username: "test1"},
			want:   url.Values{"username": []string{"test1"}},
		},
		{
			name:    "No username",
			wantErr: true,
		},
		{
			name:   "Username only",
			fields: fields{Username: "test1", FullName: "Test One", Email: "test@example.com", FirstName: "Test", LastName: "One"},
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
				Username:    tt.fields.Username,
				FullName:    tt.fields.FullName,
				Email:       tt.fields.Email,
				FirstName:   tt.fields.FirstName,
				LastName:    tt.fields.LastName,
				LDAP:        tt.fields.LDAP,
				Duo:         tt.fields.Duo,
				NeedsUpdate: tt.fields.NeedsUpdate,
			}
			got, err := u.urlValues()
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
