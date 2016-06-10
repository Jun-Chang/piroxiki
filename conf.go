package piroxiki

import "github.com/BurntSushi/toml"

type Conf map[string]struct {
	In  InConf
	Out OutConf
}

type InConf struct {
	Mail   *Mail
	Http   *Http
	Policy Policy
	Filter Filter
}

type OutConf struct {
	Mail   *Mail
	Http   *Http
	Policy Policy
}

func LoadConf(path string) (Conf, error) {
	cnf := Conf{}
	_, err := toml.DecodeFile(path, cnf)

	return cnf, err
}
