package models

import "github.com/sirupsen/logrus"

func ErrCheckf(msg string, err error) {
	if err != nil {
		logrus.Fatalf("%s %s", msg, err)
	}
}
