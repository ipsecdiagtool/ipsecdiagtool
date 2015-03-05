package ch.hsr.ipsecdiagtool.headers;

import org.jnetpcap.packet.JHeader;
import org.jnetpcap.packet.JPacket;
import org.jnetpcap.packet.JRegistry;
import org.jnetpcap.packet.RegistryHeaderErrors;
import org.jnetpcap.packet.annotate.Bind;
import org.jnetpcap.packet.annotate.Field;
import org.jnetpcap.packet.annotate.Header;
import org.jnetpcap.protocol.network.Ip4;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/*
 * ESP Implementation, capable of returning SPI and Sequence number.
 * http://www.ietf.org/rfc/rfc4303.txt
 */
@Header(length=32)
public class ESPHeader extends JHeader {

    static {
        try {
            JRegistry.register(ESPHeader.class);
        } catch (RegistryHeaderErrors e) {
            e.printStackTrace();
        }
    }

    /*
     * The ESP header is inserted after the IP header and before the next
     * layer protocol header (transport mode) or before an encapsulated IP
     * header (tunnel mode).  These modes are described in more detail below.
     */
    @Bind(to = Ip4.class)
    public static boolean bindToIp4(JPacket packet, Ip4 ip) {
        return ip.type() == 0x32; // 32 = ESP protocol
    }

    /*
     * The SPI is an arbitrary 32-bit value that is used by a receiver to
     * identify the SA to which an incoming packet is bound.  The SPI field
     * is mandatory.
     */
    @Field(offset = 136, length = 16, description = "Security Parameter Index")
    public int SPI() {
        byte[] spi = super.getByteArray(0, 4);
        //System.out.println(bytesToHex(spi));
        return byteArrayToInt(spi);
    }

    /*
     * This unsigned 32-bit field contains a counter value that increases by
     * one for each packet sent, i.e., a per-SA packet sequence number.
     */
    @Field(offset = 152, length = 16, description = "Sequence")
    public int Sequence() {
        byte[] seq = super.getByteArray(4, 4);
        //System.out.println(bytesToHex(seq));
        return byteArrayToInt(seq);
    }

    private int byteArrayToInt(byte[] b) {
        final ByteBuffer bb = ByteBuffer.wrap(b);
        bb.order(ByteOrder.BIG_ENDIAN);
        return bb.getInt();
    }

    //TODO: perhaps remove this later or put it in a different class
    //For testing hex output -- will need to be refactored to a proper place
    final protected static char[] hexArray = "0123456789ABCDEF".toCharArray();
    private String bytesToHex(byte[] bytes) {
        char[] hexChars = new char[bytes.length * 2];
        for ( int j = 0; j < bytes.length; j++ ) {
            int v = bytes[j] & 0xFF;
            hexChars[j * 2] = hexArray[v >>> 4];
            hexChars[j * 2 + 1] = hexArray[v & 0x0F];
        }
        return new String(hexChars);
    }
}
