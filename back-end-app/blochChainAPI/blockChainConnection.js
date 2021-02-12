/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('./AppUtil.js');
const { Console } = require('console');

const {channelName, chaincodeName, 
	mspOrg1, walletPath, org1UserId} = require("./fabricConstants.js");


function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}


// getConnection on specific channel with user
async function getConnection(wallet, channelName, chaincodeName, orgUserId){
	let contract;
	const ccp = buildCCPOrg1();
	const gateway = new Gateway();
	try {
		await gateway.connect(ccp, {
			wallet,
			identity: orgUserId,
			discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
		});
		const network = await gateway.getNetwork(channelName);
		contract = network.getContract(chaincodeName);
	}  finally {
		return {
			"contract": contract,
			"gateway": gateway
		};
	// Disconnect from the gateway when the application is closing
	// This will close all connections to the network
	}
}

function closeConnection(gateway) {
	gateway.disconnect();
}


async function blockChainInit() {
	try {
		// build an in memory object with the network configuration (also known as a connection profile)
		const ccp = buildCCPOrg1();
		const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
		const wallet = await buildWallet(Wallets, walletPath);
		await enrollAdmin(caClient, wallet, mspOrg1);
		await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');
		const gateway = new Gateway();
		try {
			await gateway.connect(ccp, {
				wallet,
				identity: org1UserId,
				discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
			});
			const network = await gateway.getNetwork(channelName);
			const contract = network.getContract(chaincodeName);
			console.log('\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger');
			await contract.submitTransaction('InitLedger');
			console.log('*** Result: committed');
			console.log('\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger');
			let result = await contract.evaluateTransaction('GetAllAssets');
			console.log(`*** Result: ${prettyJSONString(result.toString())}`);

		} finally {
			gateway.disconnect();
			return wallet;
		}
	} catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
		return null;
	}
}

module.exports = {
	"blockChainInit": blockChainInit,
	"getConection": getConnection,
	"closeConnetion": closeConnection,
	"prettyJSONString": prettyJSONString
};
