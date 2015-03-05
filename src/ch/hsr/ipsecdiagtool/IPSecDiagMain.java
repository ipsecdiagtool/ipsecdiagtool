package ch.hsr.ipsecdiagtool;

public class IPSecDiagMain {

    public static void main(String [] args) {
        PacketAnalyzer analyzer = new PacketAnalyzer();
        analyzer.analyzePcapFile("/home/parallels/Desktop/capture.pcap");
    }

}
