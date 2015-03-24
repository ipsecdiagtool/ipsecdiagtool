# IPSecDiagTool

[![GoDoc](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool?status.svg)](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool)
[![Build Status](https://drone.io/github.com/IPSecDiagTool/IPSecDiagTool/status.png)](https://drone.io/github.com/IPSecDiagTool/IPSecDiagTool/latest)

**YouTrack-Server:** http://sinv-56053.edu.hsr.ch

##Project structure

    YOUR_GO_WORKSPACE
        bin/
            executables                          # command executable
        pkg/
            packages                             # package object
        src/
            github.com/ipsecdiagtool/ipsecdiagtool/
                .git/                            # Git repository metadata
                main.go                          # main source file
                mtu/
                    analyze.go                   # MTU Analyzer Code
                packetloss/
                    detect.go                    # Packetloss Detector Code

