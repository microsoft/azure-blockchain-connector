using System;
using System.Threading.Tasks;
using Nethereum.Web3;

namespace nethereum_sample
{
    class Program
    {
        static Web3 web3;

        static async Task Main(string[] args)
        {
            web3 = new Web3(new AbcRpcClient(new AbcConfig
            {
                Remote = "node1-abcdemo.blockchain.azure.com:3200",
                TenantId = "microsoft.onmicrosoft.com",
                Method = AbcMethods.AadDevice
            }));
            var balance = await web3.Eth.GetBalance.SendRequestAsync("0x8F3521F692DC20A49cC4222a0ca5e010F2476AB5");
            var etherAmount = Web3.Convert.FromWei(balance.Value);
            Console.Out.WriteLine(etherAmount);
        }
    }
}