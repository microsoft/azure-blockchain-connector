namespace nethereum_sample
{
    using System;
    using System.Net.Http.Headers;
    using System.Threading.Tasks;
    using Microsoft.IdentityModel.Clients.ActiveDirectory;
    using Nethereum.JsonRpc.Client;
    using Nethereum.Web3;

    class ActiveDirectoryConfig
    {
        public const string AuthorityHost = "https://login.microsoftonline.com";
        public const string Resource = "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca"; // Resource Id for Azure Blockchain Service
        public string Tenant = "<tenant>";
        public string ClientId = "<client_id>";
        public string ClientSecret = "<client_secret>";
        public string RedirectUri = "<redirect_uri>";
    }

    class Program
    {
        static readonly ActiveDirectoryConfig config = new ActiveDirectoryConfig();
        static readonly string nodeUri = "https://<blockchain_member_name>.blockchain.azure.com:3200";

        static Web3 web3 = null;

        static async Task Main(string[] args)
        {
            var method = "aadauthcode";

            var tok = await RetrieveToken(method, config);
            PrintToken(tok);
            web3 = CreateWeb3FromOAuth2(nodeUri, tok.AccessToken);

        }

        static Web3 CreateWeb3FromOAuth2(string nodeUri, string accessToken)
        {
            var authHeader = new AuthenticationHeaderValue("Bearer", accessToken);
            return new Web3(new RpcClient(new Uri(nodeUri), authHeader, null, null));
        }

        static void PrintToken(AuthenticationResult tok)
        {
            Console.WriteLine("Access Token: " + tok.AccessToken);
        }

        static async Task<AuthenticationResult> RetrieveToken(string method, ActiveDirectoryConfig config)
        {
            var ctx = new AuthenticationContext(ActiveDirectoryConfig.AuthorityHost + "/" + config.Tenant);
            AuthenticationResult result;

            try
            {
                result = await ctx.AcquireTokenSilentAsync(ActiveDirectoryConfig.Resource, config.ClientId);
            }
            catch (Exception)
            {
                switch (method)
                {
                    case "aadauthcode":
                        result = await ctx.AcquireTokenAsync(ActiveDirectoryConfig.Resource, config.ClientId,
                            new Uri(config.RedirectUri),
                            new PlatformParameters(PromptBehavior.SelectAccount));
                        break;
                    case "aaddevice":
                        var codeResult = await ctx.AcquireDeviceCodeAsync(ActiveDirectoryConfig.Resource, config.ClientId);
                        Console.WriteLine("Open: " + codeResult.VerificationUrl);
                        Console.WriteLine("Enter: " + codeResult.UserCode);
                        result = await ctx.AcquireTokenByDeviceCodeAsync(codeResult);
                        break;
                    case "aadclient":
                        var clientCredential = new ClientCredential(config.ClientId, config.ClientSecret);
                        result = await ctx.AcquireTokenAsync(ActiveDirectoryConfig.Resource, clientCredential);
                        break;
                    default:
                        Console.WriteLine("Method not found");
                        throw new Exception("Method not found");
                }
            }

            return result;
        }
    }
}