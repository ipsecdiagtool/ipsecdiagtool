package packetloss


import(
    "log/syslog"
)

var logwriter *syslog.Writer

func InitLoger(){

    // Configure logger to write to the syslog. You could do this in init(), too.
    var e error
    logwriter, e = syslog.New(syslog.LOG_NOTICE, "IpsecDiagTool")
    if e == nil {
        logwriter.Info("IpsecDiagTool started!")
    }      
}

func InfoLog(info string){
	logwriter.Info(info)
}

func AlertLog(alert string){
	logwriter.Alert(alert)
}