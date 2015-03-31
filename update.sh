#!/usr/bin/env bash

echo "IPSecDiagTool - Updater Utility"
echo "==============================="
echo ""
sudo su
wget -N http://152.96.56.53:40000/job/IPSecDiagTool%20-%20Application/lastSuccessfulBuild/artifact/bin/ipsecdiagtool
chmod +x ipsecdiagtool
echo ""
echo "All done. Have a nice day!"