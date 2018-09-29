# Nethereum Connection Sample

This is the sample code to show how we can let Nethereum to connect to our transaction nodes directly, without the blockchain connector.

To start quickly, you can just have a look on **Nethereum/Program.cs**. Note that the region "web3 initialization" is the most important part!

If you want to have a deep look on this sample code, please see the sections below for more infomation.

# Dependency

This sample is a .Net Core 2.0 console app. It has 3 dependencies from Nuget:

    - package "Nethereum.Portable" Version="2.5.1"

    - package "Common.Logging.Core" Version="3.4.1"

    - package "Newtonsoft.Json" Version="11.0.2"

Note:
 - Here the last two packages are dependencies of "Nethereum.Portable" and haven't beed directly used in the sample code.

# How to Run the Sample Code

1. Set the parameters: 
 - Go to the region "parameters", and set the parameters following the comments.
 - If you don't have an account, see the section "New Account" for help.

2. Unlock the account:
 - As Parity and Quorum have different interface of jsonrpc method "personal_unlockAccount", we need to switch the unlock code manually.
 - For Quorum, please comment the code in the region "Parity unlock", and uncomment the code in the region "Quorum unlock".
 - For Parity, please do the opposite.

3. Build and run the code! If every thing goes right, you'll see the logs like below:
```bash
Account unlock successed!

Contract deploying finished:
    ContractAddress: 0xf8b90b91ea42954f620d159cdf157847037b72e5
    BlockHash: 0xc5dc4a030a3d4e982c92218d13e11a36f1107dd4838b6a507234d6546e62af93
    TransactionHash: 0x9d2ceeea587e99685bb0edb7a43a507fee67da3c365bc89493aa739f30b24633
    Blocknumber: 14307

Sample code finished successfully!
```

Notice: If you meet the problem that the cert of the transaction node cannot be verified, see section "Cert Verification" for help.

# Explanation of the Code

1. The most important part of the code is the region "web3 init". After init the Web3 instance in this way, you can connect to the transaction node with the Web3 instance.

2. Difference about construct Web3 with or without a managed account:
 - If a Web3 instance is constructed without a managed account, the jsonrpc method "eth_sendtransaction" is called when using this Web3 instance to send transaction.
 - If a Web3 instance is constructed with a managed account, the jsonrpc method "personal_sendtransaction" is called when using this Web3 instance to send transaction.

# New Account

Enabling Nethereum to connect to the transaction node doesn't need a blockchain account, but our sample code does (because we deploy a contract in the code). If you didn't have an account before, you can apply for an account in our sample code by doing the following steps:

1. Uncomment the code in the region "new account".

2. In the region "parameters", set the variable accountPassphase. It will become the passphase of the new account.

3. In the region "parameters", skip the setting of the variable account, but do not comment it. (For Quorum and Parity, when you are applying for an account, you only need to give it the passphase, and then the account address will be given by Quorum and Parity).

Note:
 - The new account also need to be unlocked, so do not forget to do the 2nd step in section "How to Run the Sample Code".
 - It's not recommend to apply for a new account everytime you run the sample code (or other code that need an account).

# Cert Verification

If the cert of the transaction node cannot be verified, you can simply uncomment the code in the region "clietHandler init" to tackle the issue. Here is some notice about it:

1. From some experiments, we found the reason why the cert cannot be verified is that the client doesn't know the root ca (while the cert chain is a right cert chain). 

2. We just hard coded the right root ca's issuer and thumbPrint in clientHandler.ServerCertificateCustomValidationCallback(), and verify the node's root ca.