package blockchain;

import blockchain.types.Address;

public class Special {
    public native static Address msgSender();
    public native static long msgValue();
    public native static Address txOrigin();
    public native static long gasPrice();
    public native static Address thisAddr();
    public native static long now();

    public native static void revert();
    
    // although solidity's gasLeft() returns uint256, in geth gas is uint64 and Java doc says OK for long to represent unsigned long 
    public native static long gasLeft();
}
