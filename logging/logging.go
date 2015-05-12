package logging

import (
	"fmt"
	"log/syslog"
)

type Logging struct {
	logwriter    *syslog.Writer
	alertcounter int
	alerttime    int
}

var logging Logging

// Configure logger to write to the syslog.
func InitLoger(SyslogServer string, AlertCounter int, AlertTime int) {

	//logwriter, e = syslog.New(syslog.LOG_NOTICE, "IpsecDiagTool")
	logwriter, e := syslog.Dial("udp", SyslogServer, syslog.LOG_ERR, "IpsecDiagTool")
	if e == nil && logwriter != nil {
		logwriter.Info("IpsecDiagTool started!")
	}

	logging = Logging{logwriter, AlertCounter, AlertTime}
}

func InfoLog(info string) {
	logging.logwriter.Info(info)
	fmt.Println("Syslog Info: " + info)
}

func AlertLog(alert string) {
	logging.logwriter.Alert(alert)
	fmt.Println("Syslog Alert: " + alert)
}

func AlertTime() int {
	return logging.alerttime
}

func AlertCounter() int {
	return logging.alertcounter
}
