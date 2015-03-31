#!/usr/bin/env bash

echo "IPSecDiagTool - Updater Utility"
echo "==============================="
echo ""
sudo wget -N http://152.96.56.53:40000/job/IPSecDiagTool%20-%20Application/lastSuccessfulBuild/artifact/bin/ipsecdiagtool
sudo chmod +x ipsecdiagtool
sudo wget -N http://152.96.56.53:40000/job/IPSecDiagTool%20-%20Documentation/lastSuccessfulBuild/artifact/IPSecDiagTool.pdf
echo "All done. Have a nice day!"