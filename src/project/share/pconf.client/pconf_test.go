package pconfclient

import (
	"testing"
	"strconv"
)

type SampleConfDef struct {
	foo         string  `pconf:"foo"`
	bar         int     `pconf:"bar"`
}
func (scd *SampleConfDef) Foo() string {
	return scd.foo
}
func (scd *SampleConfDef) SetFoo(v string) error {
	scd.foo = v
	return nil
}
func (scd *SampleConfDef) OnUpdateFoo(val, oldVal string) {
	return
}
func (scd *SampleConfDef) Bar() int {
	return scd.bar
}
func (scd *SampleConfDef) SetBar(v string) error {
	ival, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	scd.bar = ival
	return nil
}
func (scd *SampleConfDef) OnUpdateBar(val, oldVal string) {
	return
}

func TestRegisterConfDef(t *testing.T) {
	err := RegisterConfDef(&SampleConfDef{})
	if err != nil {
		t.Error(err)
	}
	t.Log("test normal ok.")
}


type SampleConfDef2 struct {
	hello         string    `pconf:"hello"`
	world         int
}
func (scd2 *SampleConfDef2) Hello() string {
	return scd2.hello
}
func (scd2 *SampleConfDef2) SetHello(v string) error {
	scd2.hello = v
	return nil
}
func (scd2 *SampleConfDef2) OnUpdateHello(val, oldVal string) {
	return
}

func TestRegisterConfDef2(t *testing.T) {
	err := RegisterConfDef(new(SampleConfDef2))
	if err != nil {
		t.Error(err)
	}
	t.Log("test normal ok.")
}