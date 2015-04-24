package packetloss


import(
    "log/syslog"
    "time"
    "fmt"
)

var logwriter *syslog.Writer
var alertcounter int
var alerttime int

// Configure logger to write to the syslog.
func InitLoger(SyslogServer string, AlertCounter int, AlertTime int){
	
	alertcounter = AlertCounter
	alerttime = AlertTime
	
    var e error
    //logwriter, e = syslog.New(syslog.LOG_NOTICE, "IpsecDiagTool")
    fmt.Println("Logserver: ",SyslogServer)
    logwriter, e = syslog.Dial("udp", SyslogServer, syslog.LOG_ERR, "IpsecDiagTool")
    fmt.Println(logwriter)
    if e == nil && logwriter != nil{
        logwriter.Info("IpsecDiagTool started!")
    }    
}

func InfoLog(info string){
	logwriter.Info(info)
	fmt.Println("Syslog Info")
}

func AlertLog(alert string){
	logwriter.Alert(alert)
	fmt.Println("Syslog Alert")
}

//Checks the current espmap if alert logging is necessary
func CheckLog(lostpackets []LostPacket)bool{
	var counter int
	currenttime := time.Now().Local()
	for _, v := range lostpackets {
		seconds := currenttime.Sub(v.timestamp).Seconds()		
		if(seconds < float64(alerttime)){
			counter ++
		}
	}
	return counter > alertcounter
}