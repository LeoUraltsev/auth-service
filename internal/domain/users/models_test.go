package users

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	validPwd, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				name:  Name("user"),
				email: Email{value: "validemail@gmail.com"},
				password: Password{
					validPwd,
				},
			},
			wantErr: false,
		},
		{
			name: "failure name",
			args: args{
				name:  Name(""),
				email: Email{value: "validemail@gmail.com"},
				password: Password{
					validPwd,
				},
			},
			wantErr: true,
		},
		{
			name: "failure email",
			args: args{
				name:  Name(""),
				email: Email{value: ""},
				password: Password{
					validPwd,
				},
			},
			wantErr: true,
		},
		{
			name: "failure password",
			args: args{
				name:  Name(""),
				email: Email{value: "validemail@gmail.com"},
				password: Password{
					nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	type fields struct {
		value Email
	}
	successEmailString := "success@gmail.com"
	successEmail, _ := NewEmail(successEmailString)
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				successEmail,
			},
			want: successEmailString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.value.String(); got != tt.want {
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
		{
			name: "success",
			fields: fields{
				"goodsemail@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "empty value",
			fields: fields{
				"",
			},
			wantErr: true,
		},
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
		{
			name: "success",
			n:    Name("user"),
			want: "user",
		},
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
		{
			name:    "success",
			n:       Name("user"),
			wantErr: false,
		},
		{
			name:    "empty name",
			n:       Name(""),
			wantErr: true,
		},
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
		{
			name: "success",
			args: args{
				email: "successemail@gmail.com",
			},
			want: Email{
				value: "successemail@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			args: args{
				email: "",
			},
			want:    Email{},
			wantErr: true,
		},
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
		{
			name: "success",
			args: args{
				name: "name",
			},
			want:    Name("name"),
			wantErr: false,
		},
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
		{
			name: "success",
			args: args{
				passwordHash: []byte("password"),
			},
			want:    Password{[]byte("password")},
			wantErr: false,
		},
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

	timeNow := time.Now()
	uuidNew := uuid.New()
	emailNew := Email{value: "successemail@gmail.com"}
	passwordNew := Password{[]byte("password")}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				id:           uuidNew,
				name:         "Leonard",
				email:        emailNew,
				passwordHash: passwordNew,
				isActive:     true,
				createdAt:    timeNow,
				updatedAt:    timeNow,
			},
			want: &User{
				id:           uuidNew,
				name:         "Leonard",
				email:        emailNew,
				passwordHash: passwordNew,
				isActive:     true,
				createdAt:    timeNow,
				updatedAt:    timeNow,
			},
			wantErr: false,
		},
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

func TestPassword_validate(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				hash: []byte("password"),
			},
			wantErr: false,
		},
		{
			name: "nil hash",
			fields: fields{
				hash: nil,
			},
			wantErr: true,
		},
		{
			name: "empty hash",
			fields: fields{
				hash: []byte(""),
			},
			wantErr: true,
		},
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
	successCreatedTime := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "success",
			fields: fields{
				id:           uuid.New(),
				name:         "Leonard",
				email:        Email{value: "success@gmail.com"},
				passwordHash: Password{[]byte("password")},
				isActive:     true,
				createdAt:    successCreatedTime,
				updatedAt:    time.Now().UTC(),
			},
			want: successCreatedTime,
		},
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
		{
			name: "success",
			fields: fields{
				email: Email{"success@email.com"},
			},
			want: Email{value: "success@email.com"},
		},
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
	successUserID := uuid.New()
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "success",
			fields: fields{
				id: successUserID,
			},
			want: successUserID,
		},
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
		{
			name: "success",
			fields: fields{
				isActive: true,
			},
			want: true,
		},
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
		{
			name: "success",
			fields: fields{
				name: "Leonard",
			},
			want: Name("Leonard"),
		},
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
		{
			name: "success",
			fields: fields{
				passwordHash: Password{[]byte("password")},
			},
			want: Password{[]byte("password")},
		},
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
		{
			name: "success",
			fields: fields{
				id:   uuid.New(),
				name: "Leonard",
				email: Email{
					value: "successEmail@gmail.com",
				},
				passwordHash: Password{
					[]byte("password"),
				},
				isActive:  true,
				createdAt: time.Now(),
				updatedAt: time.Now(),
			},
			args: args{
				email: Email{
					value: "NewEmail@gmail.com",
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				id:   uuid.New(),
				name: "Leonard",
				email: Email{
					value: "successEmail@gmail.com",
				},
				passwordHash: Password{
					[]byte("password"),
				},
				isActive:  true,
				createdAt: time.Now(),
				updatedAt: time.Now(),
			},
			args: args{
				email: Email{""},
			},
			wantErr: true,
		},
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
			if u.Email() == tt.fields.email && tt.wantErr == false {
				t.Errorf("UpdateEmail() Email = %v, want %v", u.Email(), tt.fields.email)
			}
			if u.updatedAt == tt.fields.updatedAt && tt.wantErr == false {
				t.Errorf("UpdateEmail() UpdatedAt = %v, want %v", u.UpdatedAt(), tt.fields.updatedAt)
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
		{
			name: "success",
			fields: fields{
				id:   uuid.New(),
				name: "Leonard",
				email: Email{
					value: "successEmail@gmail.com",
				},
				passwordHash: Password{
					[]byte("password"),
				},
				isActive:  true,
				createdAt: time.Now(),
				updatedAt: time.Now(),
			},
			args: args{
				password: Password{
					[]byte("password"),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				id:   uuid.New(),
				name: "Leonard",
				email: Email{
					value: "successEmail@gmail.com",
				},
				passwordHash: Password{
					[]byte("password"),
				},
				isActive:  true,
				createdAt: time.Now(),
				updatedAt: time.Now(),
			},
			args: args{
				password: Password{
					[]byte(""),
				},
			},
			wantErr: true,
		},
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

	successUpdatedAt := time.Now().UTC()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "success",
			fields: fields{
				id:           uuid.New(),
				name:         "Leonard",
				email:        Email{value: "success@gmail.com"},
				passwordHash: Password{[]byte("password")},
				isActive:     true,
				createdAt:    time.Now().UTC(),
				updatedAt:    successUpdatedAt,
			},
			want: successUpdatedAt,
		},
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
