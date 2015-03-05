package ch.hsr.ipsecdiagtool;
        import ch.hsr.ipsecdiagtool.headers.ESPHeader;
        import org.jnetpcap.Pcap;
        import org.jnetpcap.packet.JRegistry;
        import org.jnetpcap.packet.PcapPacket;
        import org.jnetpcap.packet.PcapPacketHandler;
        import org.jnetpcap.protocol.lan.Ethernet;
        import org.jnetpcap.protocol.network.Ip4;
        import java.util.Date;

public class IPSecDiagMain {

    public static void main(String [] args) {

        PacketAnalyzer analyzer = new PacketAnalyzer();
        analyzer.readfile("/home/parallels/Desktop/capture.pcap");
    }

}
