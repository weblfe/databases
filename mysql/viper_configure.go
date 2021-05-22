package mysql

import (
	"github.com/spf13/viper"
)

type databases struct {
	Databases DbManager `json:"databases"`
}

func NewDbMangerByViper(viper viper.Viper) (*DbManager, error) {
	var mgr = &databases{}
	if err := viper.Unmarshal(mgr); err != nil {
		return nil, err
	}
	return &mgr.Databases, nil
}
