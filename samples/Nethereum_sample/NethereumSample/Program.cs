using System;
using System.Text;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Diagnostics;
using System.Security.Cryptography.X509Certificates;
using System.Net.Security;

using Nethereum.Web3;
using Nethereum.JsonRpc.Client;
using Nethereum.Web3.Accounts.Managed;
using Nethereum.Hex.HexTypes;
namespace NethereumSample
{
    class Program
    {
        //Please set your own parameters here!!!
        //If you haven't got an account, you can uncomment the code in the region "new account",
        //set the passphase of your future account in the variable accountPassphase, and skip the setting of the variable "account".
        #region parameters        
        static string username = "<username>";
        static string password = "<password>";
        static string nodeUri = "<node_uri>";
        static string account = "<account>";
        static string accountPassphase = "<account_passphase>";        
        #endregion

        static Web3 web3 = null;
        static HttpClientHandler clientHandler = null;

        static void Main(string[] args)
        {
            //uncomment the code in this region if you met the issue that the cert of the transaction node cannot be verified.
            #region clietHandler init
            //InitHandler();
            #endregion

            //This region is the key part to let Nethereum to connect to the transaction nodes directly without the blockchain connector!
            //Following the code in this region, you can use the instance web3 to connect to the transaction nodes.            
            #region web3 init

            var authValue = Convert.ToBase64String(Encoding.UTF8.GetBytes(string.Format("{0}:{1}", username, password)));
            var authHeader = new AuthenticationHeaderValue("Basic", authValue);
            web3 = new Web3(new RpcClient(new Uri(nodeUri), authHeader, null, clientHandler));
            //Notice: clientHandler is null if the code in region "clietHandler init" remains commented.

            //if you want to construct a Web3 instance with a managed account, you can use the code below:
            //web3 = new Web3(new ManagedAccount(account, accountPassphase), new RpcClient(new Uri(nodeUri), authHeader, null, clientHandler));

            #endregion

            #region new account
            //newAccount();
            #endregion

            #region account preparation
            unlockAccount();
            #endregion

            #region test Transaction
            testTransaction();
            #endregion

            Console.ReadKey();
        }
        static void unlockAccount()
        {
            //As Parity and Quorum have different interface on jsonrpc method "personal_unlockAccount", we need to switch the unlock code manually.
            //Please make sure that, in the two regions below, only one region is uncommented.

            #region Parity unlock
            var success = web3.Personal.UnlockAccount.SendRequestAsync(account, accountPassphase, new HexBigInteger(60)).Result;
            #endregion

            #region Quorum unlock
            //var success = web3.Personal.UnlockAccount.SendRequestAsync(account, accountPassphase, 60).Result;
            #endregion

            Trace.Assert(success, "Unlock account failed! Please check your account address and the passphase!");
            Console.WriteLine("Account unlock successed!\n");
        }
        static void testTransaction()
        {
            const string contractByteCode = "0x608060405234801561001057600080fd5b50602a60008190555060df806100276000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680636d4ce63c14604e578063e5c19b2d146076575b600080fd5b348015605957600080fd5b50606060a0565b6040518082815260200191505060405180910390f35b348015608157600080fd5b50609e6004803603810190808035906020019092919050505060a9565b005b60008054905090565b80600081905550505600a165627a7a723058200e577c111b0ee4c2177cd4431abe395d21431e594a9441e820442f4ddbbe484f0029";
            const string abi = @"[{""constant"":true, ""inputs"":[],""name"":""get"",""outputs"":[{""name"":"""",""type"":""int256""}],""payable"":false,""stateMutability"":""view"",""type"":""function""},{""constant"":false,""inputs"":[{""name"":""_a"",""type"":""int256""}],""name"":""set"",""outputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""function""},{""inputs"":[],""payable"":false,""stateMutability"":""nonpayable"",""type"":""constructor""}]";
            var receiptDeploy = web3.Eth.DeployContract.SendRequestAndWaitForReceiptAsync(abi, contractByteCode, account, new HexBigInteger(4700000), new HexBigInteger(0), new HexBigInteger(0), null).Result;

            Trace.Assert(receiptDeploy.Status.Value.IsOne, "Contract Deployment failed. Check if it ran out of gas.");

            Console.WriteLine("Contract deploying finished:");
            Console.WriteLine("    ContractAddress: " + receiptDeploy.ContractAddress);
            Console.WriteLine("    BlockHash: " + receiptDeploy.BlockHash);
            Console.WriteLine("    TransactionHash: " + receiptDeploy.TransactionHash);
            Console.WriteLine("    Blocknumber: " + receiptDeploy.BlockNumber.Value);
            Console.WriteLine("\nSample code finished successfully!");

        }
        static void newAccount()
        {
            account = web3.Personal.NewAccount.SendRequestAsync(accountPassphase).Result;
            Console.WriteLine("New account successed! The account address is:");
            Console.WriteLine(account);
            Console.WriteLine();
        }
        static void InitHandler()
        {
            //We just hard coded the issuer and the thumbPrint of the Root CA in the code, 
            //and we just want to verify if the Root CA of the transaction node is the hard-encoded ca.
            const string issuer = "CN=Baltimore CyberTrust Root, OU=CyberTrust, O=Baltimore, C=IE";
            const string thumbPrint = "D4DE20D05E66FC53FE1A50882C78DB2852CAE474";
            clientHandler = new HttpClientHandler();            
            clientHandler.ServerCertificateCustomValidationCallback = (sender, cert, chain, errors) =>
            {
                if (chain == null)
                    return false;
                if (errors != SslPolicyErrors.None)
                    return false;
                X509Certificate2 cert2 = null;

                //the last element of chain.ChainElements should be the root ca
                foreach (var x in chain.ChainElements)
                {
                    cert2 = x.Certificate;
                }
                if (issuer == cert2.Issuer && thumbPrint == cert2.Thumbprint)
                    return true;
                return false;
            };
        }
    }
}
