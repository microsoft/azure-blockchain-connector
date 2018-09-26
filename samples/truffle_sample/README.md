# Truffle Connection Sample

This is the sample code to show how we can let truffle to connect to our transaction nodes directly, without the blockchain connector.

To start quickly, you can only have a look on **truffle.js**. Essentially, this sample only changes the code in **truffle.js** after
```bash
truffle init
```

If you want to have a deep look on this sample code, please see the sections below for more infomation.

# Dependency Installation

1. Install truffle:
```bash
npm install -g truffle
```

2. Install web3 package:
```bash
npm install web3@0.20.7
```

Notice:
 - only web3@0.20.1 to web3@0.20.7 can make the sample works.
 - for versions that are not higher than 0.20.0 , web3.providers.HttpProvider cannot add the authorization header to the request correctly.
 - for versions that are not lower than 1.0.0-beta.1, truffle cannot run correctly (the most possible reason is interface changing)

# Network Setting and Account Setting

1. Open **truffle.js** and change the parameters \<node_uri\>, \<username\>, \<password\> ,\<account\> and \<account_passphase\> to your own parameters.

Notice:
 - The "Account Unlock" part is not always needed. We add it here because our sample will deploy a contract, and this action needs an account.
 - If you didn't have any accounts, just comment the code in "Account Unlock" and "Sender Account Assignation" part and continue. Later you'll know how to get an account.
 - In your own application, you may choose to unlock an account at some other position, or even you may not need to unlock any account.

2. In the folder, run:
```bash
truffle console --network net1
```
then you can see "truffle(net1)>" in the new line, which indicates you've entered truffle's console. 

Later when you want to exit the console, you can run:
```bash
.exit
```

If you met the connection error "Could not connect to your Ethereum client. Please check that ......", Please see the section "SSL Cert Injection" for help.

Notice:
    - net1 is the name of the network definition in **truffle.js**, and you can change the name by yourself.

3. Apply for a new account (If you already have an account, you can skip this step).
In the truffle console, run:
```bash
web3.personal.newAccount("<account_passphase>")
```
and you'll see a string, which is your account address (of course, it can be filled into all the <account> metioned above).

Then you need to Uncomment the code in **truffle.js**, which you've commented in step 1.

#Truffle Migrate

1. Run:
```bash
truffle migrate --network net1 --reset
```
and you'll see logs like below:
```bash
Using network 'net3'.

Running migration: 1_initial_migration.js
  Replacing Migrations...
  ... 0xaad543c7e83c295727cc1c323a5de53a7e55a3b1c98294ec20f8a8047b44eca1
  Migrations: 0x8054a2f9023725ae8dc5e5e2bfffd61e9d6f9a38
Saving successful migration to network...
  ... 0x263ca1bf03cd5b1d93a2e152fb326bf76edaca67e91c90101d42e7b4687fa8a7
Saving artifacts...
```
which indicates the sample runs successfully!

If you met the connection error "Could not connect to your Ethereum client. Please check that ......", Please see the section "SSL Cert Injection" for help.

#SSL Cert Injection
The connection error "Could not connect to your Ethereum client. Please check that ......" may caused by that web3 consider the ssl cert of the transaction node as an unsafe cert. If it's caused by this reason, we can recover it by injecting the cert manually:

1. Uncomment the code in "SSL Cert Injection" part.

2. Change <cert_path> to the path of the cert on your computer (How to get the cert will be updated later).

3. Install the package ssl-root-cas:
 ```bash
npm install ssl-root-cas
```