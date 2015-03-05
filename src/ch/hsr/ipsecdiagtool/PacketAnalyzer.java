package ch.hsr.ipsecdiagtool;

import ch.hsr.ipsecdiagtool.headers.ESPHeader;
import org.jnetpcap.Pcap;
import org.jnetpcap.packet.PcapPacket;
import org.jnetpcap.packet.PcapPacketHandler;

import java.util.Date;

public class PacketAnalyzer {
    final ESPHeader espHeader = new ESPHeader();
    StringBuilder errorBuffer = new StringBuilder();

    public void analyzePcapFile(String filename) {

        Pcap pcap = Pcap.openOffline(filename, errorBuffer);

        if (pcap == null) {
            System.err.printf("Error while opening device for capture: "
                    + errorBuffer.toString());
            return;
        }

        PcapPacketHandler<String> jpacketHandler = new PcapPacketHandler<String>() {
            public void nextPacket(PcapPacket packet, String user) {

                //Ignore non-ESP packets
                if(packet.hasHeader(espHeader)){
                    System.out.println("ESP-Packet from "+ new Date(packet.getCaptureHeader().timestampInMillis())
                            +" Sequence-NR."+espHeader.Sequence()+" SPI."+espHeader.SPI());
                }
            }
        };

        pcap.loop(100, jpacketHandler, null);

        //Print JRegistery (shows all registered protocols)
        //System.out.println(JRegistry.toDebugString());
        pcap.close();
    }
}
