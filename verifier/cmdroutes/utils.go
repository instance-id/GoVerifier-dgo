package cmdroutes

import "github.com/sirupsen/logrus"

func TrimInvoice(invoice string) {

}

func ErrCheckf(msg string, err error) {
	if err != nil {
		logrus.Fatalf("%s %s", msg, err)
	}
}
