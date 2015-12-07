package deploy

import (
	"errors"
	"fmt"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ini.v1"
)

type Option struct {
	GitOptions map[string]GitOption
}

func NewOption(iniFile *ini.File) (*Option, error) {
	deployItems := iniFile.Section("deploy").KeysHash()
	if len(deployItems) == 0 {
		log15.Warn("Deploy.Init.NoConf")
		return nil, nil
	}

	dOpt := &Option{}

	for _, name := range deployItems {
		section := iniFile.Section("deploy." + name)
		typeName := section.Key("type").String()
		if typeName == "" {
			log15.Error("Deploy.Init.UnknownType.[" + name + "]")
			continue
		}
		fmt.Println("deploy read type", typeName)
		switch typeName {
		case TYPE_GIT:
			opt := GitOption{
				Message: "site updated at {now}",
			}
			if err := section.MapTo(&opt); err != nil {
				return nil, err
			}
			if err := opt.isValid(); err != nil {
				return nil, err
			}
			if len(dOpt.GitOptions) == 0 {
				dOpt.GitOptions = make(map[string]GitOption)
			}
			dOpt.GitOptions[name] = opt
		default:
			return nil, errors.New("Unknown type " + typeName)
		}
	}

	if len(dOpt.GitOptions) == 0 {
		return nil, nil
	}
	return dOpt, nil
}
