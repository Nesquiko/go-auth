package db

import (
	"testing"
)

func TestConnectDBIsNil(t *testing.T) {
	if DBConnection != nil {
		t.Error("Uninitialized connection to db is not nil")
	}
}

func TestConnectDBErrors(t *testing.T) {
	type args struct {
		driver string
		dsn    string
	}

	conFailMsg := "Connection should be nil to indicate error"
	errFailMsg := "Err cannot be nil"

	tests := []struct {
		name string
		args args
	}{
		{
			name: "InvalidDSN",
			args: args{driver: "mysql", dsn: "invalidDSN"},
		},
		{
			name: "InvalidDriver",
			args: args{driver: "invalid", dsn: "root:passwd@tcp(127.0.0.1:3306)/users?"},
		},
		{
			name: "ConnectionRefused",
			args: args{driver: "mysql", dsn: "root:passwd@tcp(127.0.0.1:3306)/users?"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConnectDB(tt.args.driver, tt.args.dsn)

			if DBConnection != nil {
				t.Error(conFailMsg)
			}
			if err == nil {
				t.Error(errFailMsg)
			}
		})
	}
}

func Test_connectErrors(t *testing.T) {
	type args struct {
		driver string
		dsn    string
	}

	conFailMsg := "Connection should be nil to indicate error"
	errFailMsg := "Err cannot be nil"

	tests := []struct {
		name string
		args args
	}{
		{
			name: "InvalidDSN",
			args: args{driver: "mysql", dsn: "invalidDSN"},
		},
		{
			name: "InvalidDriver",
			args: args{driver: "invalid", dsn: "root:passwd@tcp(127.0.0.1:3306)/users?"},
		},
		{
			name: "ConnectionRefused",
			args: args{driver: "mysql", dsn: "root:passwd@tcp(127.0.0.1:3306)/users?"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := connect(tt.args.driver, tt.args.dsn)

			if c != nil {
				t.Error(conFailMsg)
			}
			if err == nil {
				t.Error(errFailMsg)
			}
		})
	}
}

func TestMySQLDSNConfig(t *testing.T) {
	user := "testUser"
	passwd := "testPasswd"
	addr := "0.0.0.0:4269"
	dbname := "database"

	cfg := MySQLDSNConfig(user, passwd, addr, dbname)

	if cfg.User != user {
		t.Errorf("Expected user to be %s, but was %s", user, cfg.User)
	}
	if cfg.Passwd != passwd {
		t.Errorf("Expected passwd to be %s, but was %s", passwd, cfg.Passwd)
	}
	if cfg.Addr != addr {
		t.Errorf("Expected addr to be %s, but was %s", addr, cfg.Addr)
	}
	if cfg.DBName != dbname {
		t.Errorf("Expected dbname to be %s, but was %s", dbname, cfg.DBName)
	}
}
