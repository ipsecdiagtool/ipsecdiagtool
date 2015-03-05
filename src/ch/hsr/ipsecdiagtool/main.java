package ch.hsr.ipsecdiagtool;

import ch.hsr.ipsecdiagtool.headers.ESPHeader;
import org.jnetpcap.Pcap;
import org.jnetpcap.packet.JRegistry;
import org.jnetpcap.packet.PcapPacket;
import org.jnetpcap.packet.PcapPacketHandler;
import org.jnetpcap.protocol.lan.Ethernet;
import org.jnetpcap.protocol.network.Ip4;
import java.util.Date;

public class Main {

    public static void main(String [] args)
    {
        System.out.println("hello world");
        readfile();
    }

    public static void readfile() {

        final ESPHeader ipsec = new ESPHeader();
        final Ip4 ipv4 = new Ip4();
        final Ethernet ethernet = new Ethernet();

        String filename = "capture.pcap";
        StringBuilder errbuf = new StringBuilder(); // For any error msgs

        Pcap pcap = Pcap.openOffline(filename, errbuf);

        if (pcap == null) {
            System.err.printf("Error while opening device for capture: "
                    + errbuf.toString());
            return;
        }

        PcapPacketHandler<String> jpacketHandler = new PcapPacketHandler<String>() {

            public void nextPacket(PcapPacket packet, String user) {

                System.out.printf("Received packet at %s caplen=%-4d len=%-4d %s\n",
                        new Date(packet.getCaptureHeader().timestampInMillis()),
                        packet.getCaptureHeader().caplen(),  // Length actually captured
                        packet.getCaptureHeader().wirelen(), // Original length
                        user                                 // User supplied object
                );

                System.out.println(
                        "Time:" + new Date(packet.getCaptureHeader().timestampInMillis())
                );


                if(packet.hasHeader(ipsec)){
                    System.out.println("hello world.");
                }
            }
        };

        pcap.loop(100, jpacketHandler, "read");
        System.out.println(JRegistry.toDebugString());
        pcap.close();
    }
}
