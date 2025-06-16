package users

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		name     Name
		email    Email
		password Password
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Email{
				value: tt.fields.value,
			}
			if got := e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail_validate(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Email{
				value: tt.fields.value,
			}
			if err := e.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestName_String(t *testing.T) {
	tests := []struct {
		name string
		n    Name
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_validate(t *testing.T) {
	tests := []struct {
		name    string
		n       Name
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    Email
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Name
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPassword(t *testing.T) {
	type args struct {
		passwordHash []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Password
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPassword(tt.args.passwordHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.name, tt.args.email, tt.args.passwordHash, tt.args.isActive, tt.args.createdAt, tt.args.updatedAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_String(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Password{
				hash: tt.fields.hash,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword_validate(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Password{
				hash: tt.fields.hash,
			}
			if err := p.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_CreatedAt(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.CreatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Email(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   Email
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.Email(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ID(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Name(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   Name
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Password(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   Password
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.Password(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	type args struct {
		email Email
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if err := u.UpdateEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	type args struct {
		password Password
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if err := u.UpdatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_UpdatedAt(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		name         Name
		email        Email
		passwordHash Password
		isActive     bool
		createdAt    time.Time
		updatedAt    time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				id:           tt.fields.id,
				name:         tt.fields.name,
				email:        tt.fields.email,
				passwordHash: tt.fields.passwordHash,
				isActive:     tt.fields.isActive,
				createdAt:    tt.fields.createdAt,
				updatedAt:    tt.fields.updatedAt,
			}
			if got := u.UpdatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
