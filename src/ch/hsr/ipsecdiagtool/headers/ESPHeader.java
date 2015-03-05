package ch.hsr.ipsecdiagtool.headers;

import org.jnetpcap.packet.JHeader;
import org.jnetpcap.packet.JPacket;
import org.jnetpcap.packet.JRegistry;
import org.jnetpcap.packet.RegistryHeaderErrors;
import org.jnetpcap.packet.annotate.Bind;
import org.jnetpcap.packet.annotate.Field;
import org.jnetpcap.packet.annotate.Header;
import org.jnetpcap.protocol.network.Ip4;

@Header(length=32)
public class ESPHeader extends JHeader {

    static {
        try {
            JRegistry.register(ESPHeader.class);
        } catch (RegistryHeaderErrors e) {
            e.printStackTrace();
        }
    }

    @Bind(to = Ip4.class)
    public static boolean bindToIp4(JPacket packet, Ip4 ip) {
        return ip.type() == 0x32; // 32 = ESP protocol
    }

    @Field(offset = 136, length = 16)
    public int SPI() {
        return super.getUByte(0);
    }

    @Field(offset = 152, length = 16)
    public int Sequence() {
        return super.getUByte(0);
    }
}
