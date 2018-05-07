package blockchain;

import blockchain.types.Address;

public class Special {
    public native static Address msgSender();
    // although solidity's gasLeft() returns uint256, in geth gas is uint64 and Java doc says OK for long to represent unsigned long 
    public native static long gasLeft();
}
