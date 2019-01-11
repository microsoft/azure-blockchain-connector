using System;
using System.Net.Http.Headers;
using System.Runtime.CompilerServices;
using System.Text;
using System.Threading.Tasks;
using Microsoft.IdentityModel.Clients.ActiveDirectory;
using Nethereum.JsonRpc.Client;
using Nethereum.Web3;
using Org.BouncyCastle.Math.EC;

namespace nethereum_sample
{
    class AADConfig
    {
        public string AuthorityHost = "https://login.microsoftonline.com";
        public string Tenant = "<tenant>";
        public string ClientId = "<client_id>";
        public string ClientSecret = "<client_secret>";
        public string Resource = "<resource>";
        public string RedirectUri = "<redirect_uri>";
    }

    class Program
    {
        static string nodeUri = "<node_uri>";

        // Basic Authentication settings
        static string username = "<username>";
        static string password = "<password>";

        // AAD OAuth2 settings
        static AADConfig config = new AADConfig();

        static Web3 web3 = null;

        static async Task Main(string[] args)
        {
            var method = "aadauthcode";

            switch (method)
            {
                case "":
                case "basic":

                    // Basic Authentication
                    web3 = web3FromBasicAuth(nodeUri, username, password);

                    break;
                case "aadauthcode":
                case "aaddevice":
                case "aadclient":

                    // AAD OAuth2
                    var tok = await retrieveToken(method, config);
                    printToken(tok);
                    web3 = web3FromOAuth2(nodeUri, tok.AccessToken);

                    break;
                default:
                    return;
            }
        }

        static Web3 web3FromBasicAuth(string nodeUri, string username, string password)
        {
            var authValue =
                Convert.ToBase64String(Encoding.UTF8.GetBytes($"{username}:{password}"));
            var authHeader = new AuthenticationHeaderValue("Basic", authValue);
            return new Web3(new RpcClient(new Uri(nodeUri), authHeader, null, null));
        }

        static Web3 web3FromOAuth2(string nodeUri, string accessToken)
        {
            var authHeader = new AuthenticationHeaderValue("Bearer", accessToken);
            return new Web3(new RpcClient(new Uri(nodeUri), authHeader, null, null));
        }

        static void printToken(AuthenticationResult tok)
        {
            Console.WriteLine("Access: " + tok.AccessToken);
        }

        static async Task<AuthenticationResult> retrieveToken(string method, AADConfig config)
        {
            // aadauthcode and aaddevice method use fixed settings
            if (method == "aadauthcode" || method == "aaddevice")
            {
                config.AuthorityHost = "https://login.microsoftonline.com";
                config.Tenant = "microsoft.onmicrosoft.com";
                config.ClientId = "a8196997-9cc1-4d8a-8966-ed763e15c7e1";
                config.ClientSecret = null;
                config.Resource = "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca";
                config.RedirectUri = "urn:ietf:wg:oauth:2.0:oob";
            }

            var authority = config.AuthorityHost + "/" + config.Tenant;

            var ctx = new AuthenticationContext(authority);
            AuthenticationResult result = null;

            try
            {
                result = await ctx.AcquireTokenSilentAsync(config.Resource, config.ClientId);
            }
            catch (Exception)
            {
                switch (method)
                {
                    case "aadauthcode":
                        result =   paramconfi
                        break;
                    case "aaddevice":
                        var codeResult = await ctx.AcquireDeviceCodeAsync(config.Resource, config.ClientId);
                        Console.WriteLine("Open: " + codeResult.VerificationUrl);
                        Console.WriteLine("Enter: " + codeResult.UserCode);
                        result = await ctx.AcquireTokenByDeviceCodeAsync(codeResult);
                        break;
                    case "aadclient":
                        var clientCredential = new ClientCredential(config.ClientId, config.ClientSecret);
                        result = await ctx.AcquireTokenAsync(config.Resource, clientCredential);
                        break;
                    default:
                        Console.WriteLine("method not found");
                        throw new Exception("method not found");
                }
            }

            return result;
        }
    }
}