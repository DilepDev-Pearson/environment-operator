package config

import (
	"fmt"
	"testing"

	"github.com/pearsontechnology/environment-operator/pkg/config"
)

func TestEnvironmentsBitesize(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		t.Run("single environment name", testSingleName)
		t.Run("number of environments match", testNumberOfEnvironments)
	})

	t.Run("invalid config", testInvalidConfig)
}

func testSingleName(t *testing.T) {
	cfg := `
  project: test
  environments:
  - name: Abr
  `
	configuration, err := config.Load(cfg)
	if err != nil {
		t.Errorf("Config: unexpected error %s", err.Error())
	}

	c := configuration.Environments[0]

	if c.Name != "Abr" {
		t.Errorf("Config: error on Name (%v)", c.Name)
	}
}

func testNumberOfEnvironments(t *testing.T) {
	cfg := `
  project: test
  environments:
  - name: first
  - name: second
  `
	expected := 2

	configuration, err := config.Load(cfg)
	if err != nil {
		t.Errorf("Config: unexpected error %s", err.Error())
	}

	if len(configuration.Environments) != expected {
		t.Errorf("Config: error on Environment count. Expected: %d, Actual: %d",
			expected,
			len(configuration.Environments),
		)
	}
}

func testInvalidConfig(t *testing.T) {
	var saTests = []struct {
		Cfg      string
		Expected string
		Cause    string
	}{
		{
			`
      project: test
      environments:
      - name: zzz
        vo:
          : nono
      `,
			"yaml: line 5: did not find expected key",
			"invalid yaml",
		},
		{
			`
      project: test
      environments:
      - services:
        - name: service one
        - name: service two
      `,
			"environment.Name: zero value",
			"missing environment name",
		},
		{
			`
      project: test
      environments:
      - name: o
        services:
        - n: 1
      `,
			"environment.service.Name: zero value",
			"missing service name",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        namespace: namespace_invalid
      `,
			"environment.Namespace: regular expression mismatch",
			"invalid namespace",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        services:
        - name: Service1
          deployment:
            method: invalid_method
      `,
			"environment.service.deployment.Method: regular expression mismatch",
			"invalid service deployment method",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        deployment:
          method: invalid_method
        services:
        - name: Service1
      `,
			"environment.deployment.Method: regular expression mismatch",
			"invalid deployment method",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        deployment:
          method: bluegreen
          mode: man
        services:
        - name: Service1
      `,
			"environment.deployment.Mode: regular expression mismatch",
			"invalid deployment mode",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        deployment:
          method: bluegreen
          mode: manual
          active: red
        services:
        - name: Service1
      `,
			"environment.deployment.Active: regular expression mismatch",
			"invalid deployment active",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        services:
          - name: Service1
            health_check:
              command: command
      `,
			"environment.service.health_check.yaml: unmarshal errors:\n  line 8: cannot unmarshal !!str `command` into []string",
			"invalid service health check",
		},
		{
			`
      project: test
      environments:
      - name: Abr
        services:
          - name: Service1
            health_check:
              cmd: command
      `,
			"environment.service.health_check: unknown fields (cmd)",
			"invalid key in service health check",
		},
	}

	for idx, tst := range saTests {
		t.Run(tst.Cause, func(t *testing.T) {
			configuration, err := config.Load(tst.Cfg)
			if err != nil {
				if err.Error() != tst.Expected {
					t.Errorf("Error in %d test:\nEXPECTED:\n%s\n--\nACTUAL:\n%s\n",
						idx,
						tst.Expected,
						err.Error(),
					)
				}
			} else {
				fmt.Printf("%+v", configuration)
				t.Errorf("Config: no error on %s", tst.Cause)
			}
		})

	}

}
