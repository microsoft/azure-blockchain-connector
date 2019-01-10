const Web3 = require('web3');

(function () {
    const username = '<username>';
    const password = '<password>';
    const nodeUri = '<node_uri>';

    const web3 = new Web3();
    let provider;

    provider = new web3.providers.HttpProvider(`http://${username}:${password}@${nodeUri}`);

    // or
    provider = new web3.providers.HttpProvider("<node_uri>", 0, "<username>", "<password>");

    web3.setProvider(provider);

    // your code

})();