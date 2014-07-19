package scan

import (
	"reflect"
	"testing"

	"github.com/sourcegraph/srclib/config"
	"github.com/sourcegraph/srclib/container"
	"github.com/sourcegraph/srclib/unit"
)

type DummyPackage struct {
	Dir string
}

func (_ DummyPackage) ID() string      { return "foo" }
func (_ DummyPackage) Name() string    { return "foo" }
func (_ DummyPackage) RootDir() string { return "foo" }
func (p DummyPackage) Paths() []string { return []string{p.Dir} }

type DummyScanner struct{}

func (_ DummyScanner) Scan(dir string, c *config.Repository) ([]unit.SourceUnit, error) {
	return []unit.SourceUnit{DummyPackage{"foo"}}, nil
}

func unregister(name string) {
	delete(Scanners, name)
}

func TestSourceUnits(t *testing.T) {
	Register("dummy", DummyScanner{})
	defer unregister("dummy")

	oldRunner := container.DefaultRunner
	container.DefaultRunner = &container.MockRunner{}
	defer func() {
		container.DefaultRunner = oldRunner
	}()

	wantUnits := []unit.SourceUnit{DummyPackage{"foo"}}

	units, err := SourceUnits("qux", &config.Repository{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(units, wantUnits) {
		t.Errorf("got units %v, want %v", units, wantUnits)
	}
}
