package ch.hsr.ipsecdiagtool;

import ch.hsr.ipsecdiagtool.headers.ESPHeader;
import org.jnetpcap.Pcap;
import org.jnetpcap.packet.JRegistry;
import org.jnetpcap.packet.PcapPacket;
import org.jnetpcap.packet.PcapPacketHandler;
import org.jnetpcap.protocol.lan.Ethernet;
import org.jnetpcap.protocol.network.Ip4;

import java.util.Date;

public class PacketAnalyzer {

    public static void readfile(String filename) {
        final ESPHeader espHeader = new ESPHeader();
        final Ip4 ipv4 = new Ip4();
        final Ethernet ethernet = new Ethernet();
        StringBuilder errbuf = new StringBuilder(); // For any error msgs

        Pcap pcap = Pcap.openOffline(filename, errbuf);

        if (pcap == null) {
            System.err.printf("Error while opening device for capture: "
                    + errbuf.toString());
            return;
        }

        PcapPacketHandler<String> jpacketHandler = new PcapPacketHandler<String>() {
            public void nextPacket(PcapPacket packet, String user) {

                //Ignore non-ESP packets
                if(packet.hasHeader(espHeader)){
                    System.out.println(
                            "ESP Packet - Time:" + new Date(packet.getCaptureHeader().timestampInMillis())
                    );
                }
            }
        };

        pcap.loop(100, jpacketHandler, "read");

        //Debug info for headers:
        System.out.println(JRegistry.toDebugString());
        pcap.close();
    }
}
