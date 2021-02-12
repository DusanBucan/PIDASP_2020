var express = require('express');
var router = express.Router();

var app = require("../app");
var {getConection, closeConnetion, prettyJSONString} = require("../blochChainAPI/blockChainConnection.js");
const { chaincodeName, channelName, org1UserId } = require('../blochChainAPI/fabricConstants.js');


/* GET home page. */
router.get('/', async (req, res, next) => {

  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  console.log('\n--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger');
	let result = await contract.evaluateTransaction('GetAllAssets');
  result = JSON.parse(result.toString());
  closeConnetion(gateway);

  res.json(result);
  // res.render('index', { title: 'Express' });
});


router.get('/:assetId', async (req, res, next) => {

  const assetId = req.params.assetId;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  result = await contract.evaluateTransaction('ReadAsset', assetId);
	result = JSON.parse(result.toString());
  closeConnetion(gateway);

  res.json(result);
});


router.post("/",  async (req, res, next) => {

  const newAsset = req.body;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  result = await contract.submitTransaction('CreateAsset', 
          newAsset.ID, newAsset.Color, newAsset.Size, newAsset.Owner, newAsset.AppraisedValue);
			console.log('*** Result: committed');
	if (`${result}` !== '') {
    result = JSON.parse(result.toString());
	} else {
    result = "";
  }

  closeConnetion(gateway);
  res.json(result);
});


module.exports = router;
