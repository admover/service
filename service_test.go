// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

package service_test

import (
	"os"
	"testing"
	"time"

	"github.com/admover/service"
)

func TestRunInterrupt(t *testing.T) {
	p := &program{}
	sc := &service.Config{
		Name: "go_service_test",
	}
	s, err := service.New(p, sc)
	if err != nil {
		t.Fatalf("New err: %s", err)
	}

	go func() {
		<-time.After(1 * time.Second)
		interruptProcess(t)
	}()

	go func() {
		for i := 0; i < 25 && p.numStopped == 0; i++ {
			<-time.After(200 * time.Millisecond)
		}
		if p.numStopped == 0 {
			t.Fatal("Run() hasn't been stopped")
		}
	}()

	if err = s.Run(); err != nil {
		t.Fatalf("Run() err: %s", err)
	}
}

const testInstallEnv = "TEST_USER_INSTALL"

// Should always run, without asking for any permission
func TestUserRunInterrupt(t *testing.T) {
	if os.Getenv(testInstallEnv) != "1" {
		t.Skipf("env %q is not set to 1", testInstallEnv)
	}
	p := &program{}
	options := make(service.KeyValue)
	options["UserService"] = true
	sc := &service.Config{
		Name:   "go_user_service_test",
		Option: options,
	}
	s, err := service.New(p, sc)
	if err != nil {
		t.Fatalf("New err: %s", err)
	}
	err = s.Install()
	if err != nil {
		t.Errorf("Install err: %s", err)
	}
	err = s.Uninstall()
	if err != nil {
		t.Fatalf("Uninstall err: %s", err)
	}
}

type program struct {
	numStopped int
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	p.numStopped++
	return nil
}
