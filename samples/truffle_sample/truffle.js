/*
 * NB: since truffle-hdwallet-provider 0.0.5 you must wrap HDWallet providers in a 
 * function when declaring them. Failure to do so will cause commands to hang. ex:
 * ```
 * mainnet: {
 *     provider: function() { 
 *       return new HDWalletProvider(mnemonic, 'https://mainnet.infura.io/<infura-key>') 
 *     },
 *     network_id: '1',
 *     gas: 4500000,
 *     gasPrice: 10000000000,
 *   },
 */


const web3 = require('web3');
// const {authCodeGrant, clientCredentialsGrant} = require("./aad");
//===========SSL Cert Injection Begin============

/*
 *const https = require('https');
 *const fs =require('fs');
 *const ssl_root_cas = require('ssl-root-cas/latest');
 *var ca = fs.readFileSync('<cert_path>');
 *https.globalAgent.options.ca=require('ssl-root-cas/latest').create();
 *https.globalAgent.options.ca=ssl_root_cas.create();
 *https.globalAgent.options.ca.push(ca)
 */

//===========SSL Cert Injection End============

module.exports = {
    networks: {
        net1: {
            provider: function () {

                let provider;

                // Basic Authentication
                provider = new web3.providers.HttpProvider("<node_uri>", 0, "<username>", "<password>");

                // AAD OAuth2
                // As the provider() should return a provider synchronously, it is not supported to doing the async oauth2 authentication here.
                // Truffle once supports loading an async provider, but now it does not. See https://github.com/trufflesuite/truffle/pull/1073
                // You may run 'node aad.js' to get tokens manually, and here demos how to use the token with the official http provider.
                // Please check the configs inside 'aad.js'.
                const accessToken = "<access_token>";
                const authHeader = {name: 'Authorization', value: 'Bearer ' + accessToken};
                provider = new web3.providers.HttpProvider("<node_uri>", 0, "", "", [authHeader]);

                //The "Account Unlock" part is not needed. We add it here because our sample will deploy a contract, and this action needs an account.
                //In your own application, you may not need to unlock any account, and also you can choose to unlock an account at some other position.

                //==============Account Unlock Begin=============

                // var web3Instance = new web3(provider);
                // web3Instance.personal.unlockAccount("<account>", "<account_passphase>");

                //==============Account Unlock End=============

                return provider;
            },
            network_id: "*",
            gas: 4500000,
            gasPrice: 0,

            //================Sender Account Assignation Begin==============

            from: "<account>"

            //================Sender Account Assignation End==============
        },

    }
};
