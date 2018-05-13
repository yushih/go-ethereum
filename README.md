# J2ChainChain: a block chain that runs Java smart contract

J2ChainChain is a mod(ification) of geth. Instead of running EVM bytecode, it runs JVM bytecode. Users of this blockchain write smart contracts in Java and execute the .class file as smart contract code.

For example, the ERC20 for J2ChainChain would look like:
~~~~
import java.util.HashMap;
import java.lang.Long;

import blockchain.types.Address;
import blockchain.Special;

class FatOtakuHappyToken {
    long _totalSupply;
    HashMap<Address, Long> balances;
    HashMap<Address, HashMap<Address, Long>> allowed;
    
    FatOtakuHappyToken (long totalSupply) {
        _totalSupply = totalSupply;

        balances = new HashMap<Address, Long>();
        balances.put(Special.msgSender(), new Long(totalSupply));

        allowed = new HashMap<Address, HashMap<Address, Long>>();
    }

    public long totalSupply() {
        return _totalSupply;
    }

    // ------------------------------------------------------------------------
    // Get the token balance for account `tokenOwner`
    // ------------------------------------------------------------------------
    public long balanceOf(Address tokenOwner) {
        return (long)balances.get(tokenOwner);
    }

    void add(HashMap<Address, Long> mapping, Address account, long tokens) {
        long newBalance = (long)mapping.get(account)+tokens;
        mapping.put(account, new Long(tokens));
    }

    // ------------------------------------------------------------------------ 
    // Transfer the balance from token owner's account to `to` account
    // - Owner's account must have sufficient balance to transfer
    // - 0 value transfers are allowed
    // ------------------------------------------------------------------------
    public boolean transfer(Address to, long tokens) {
        Address from = Special.msgSender();
        if (balanceOf(from)<tokens) {
            return false;
        } else {
            add(balances, from, -tokens);
            add(balances, to, tokens);
            return true;
        }
    }

    // ------------------------------------------------------------------------
    // Token owner can approve for `spender` to transferFrom(...) `tokens`
    // from the token owner's account
    //
    // https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20-token-standard.md
    // recommends that there are no checks for the approval double-spend attack
    // as this should be implemented in user interfaces 
    // ------------------------------------------------------------------------
    public boolean approve(Address spender, long tokens) {
        Address sender = Special.msgSender();
        if (allowed.get(sender) == null) {
            allowed.put(sender, new HashMap<Address, Long>());
        }
        allowed.get(sender).put(spender, new Long(tokens));
        return true;
    }

    // ------------------------------------------------------------------------
    // Returns the amount of tokens approved by the owner that can be
    // transferred to the spender's account
    // ------------------------------------------------------------------------
    public long allowance(Address tokenOwner, Address spender) {
        return (long)allowed.get(tokenOwner).get(spender);
    }

    // ------------------------------------------------------------------------
    // Transfer `tokens` from the `from` account to the `to` account
    // 
    // The calling account must already have sufficient tokens approve(...)-d
    // for spending from the `from` account and
    // - From account must have sufficient balance to transfer
    // - Spender must have sufficient allowance to transfer
    // - 0 value transfers are allowed
    // ------------------------------------------------------------------------
    public boolean transferFrom(Address from, Address to, long tokens) {
        Address sender = Special.msgSender();
        if (allowed.get(from).get(sender) < tokens) {
            return false;
        }
        add(allowed.get(from), sender, -tokens);
        add(balances, from, -tokens);
        add(balances, to, tokens);
        return true;
    }
}
~~~~

JVM is based on [https://github.com/zxh0/jvmgo-book/](https://github.com/zxh0/jvmgo-book/).

### Development status
J2ChainChain is mostly a hobby project. It is runnable but not yet production-ready. There are many rough edges, which I will work on very slowly.


* JRE stripping and sandboxing

* Unhandled exception handling

* Support Ethereum logs

* Computation and storage effiency measuring and optimization.

* Support deploying and using library code.

* Support default handler (as Solidity does).

* Use java.math.BigInteger for integer value.

* Allow array parameter

* Object.hashCode
