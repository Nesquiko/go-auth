package db

import (
	"testing"
)

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

func Test_connectInvalidDSN(t *testing.T) {
	c, err := connect("mysql", "invalidDSN")

	if c != nil {
		t.Error("Connection should be nil to indeicate error")
	}
	if err == nil {
		t.Error("Err cannot be nil")
	}
}

func Test_connectInvalidDriver(t *testing.T) {
	c, err := connect("invalid", "root:passwd@tcp(127.0.0.1:3306)/users?")

	if c != nil {
		t.Error("Connection should be nil to indeicate error")
	}
	if err == nil {
		t.Error("Err cannot be nil")
	}
}

func Test_connectConnectionRefused(t *testing.T) {
	c, err := connect("mysql", "root:passwd@tcp(127.0.0.1:3306)/users?")

	if c != nil {
		t.Error("Connection should be nil to indeicate error")
	}
	if err == nil {
		t.Error("Err cannot be nil")
	}
}
