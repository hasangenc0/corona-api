package environment

import (
	"github.com/Netflix/go-env"
	"log"
)

type Environment struct {
	Env    string `env:"ENV"`
	Port   string `env:"PORT"`
	Extras env.EnvSet
}

func Get() Environment {
	var environment Environment
	es, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	// Remaining environment variables.
	environment.Extras = es

	return environment
}
