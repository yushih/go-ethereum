package blockchain.types;

public class Address {
    public byte[] bytes;
    public native long balance(); 

    public native Object call(String methodName, long value, Object... args);

    @Override
    public native int hashCode();
    
    @Override
    public native boolean equals(Object obj);
}
