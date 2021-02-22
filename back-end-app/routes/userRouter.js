var express = require('express');
var router = express.Router();

var app = require("../app");
var {getConection, closeConnetion, prettyJSONString} = require("../blochChainAPI/blockChainConnection.js");
const { chaincodeName, channelName, org1UserId } = require('../blochChainAPI/fabricConstants.js');


router.get('/all', async (req, res, next) => {

    let wallet = req.app.get("fabricWallet");
    let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
    try {
        let result = await contract.evaluateTransaction('GetAllUsers');
        result = JSON.parse(result.toString());
        res.json(result);
    } catch (err) {
        res.json(err.message);
    } finally {
        closeConnetion(gateway);
    }
});

router.get('/:personId', async (req, res, next) => {
    const personId = req.params.personId;
    let wallet = req.app.get("fabricWallet");
    let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
    try {
      result = await contract.evaluateTransaction('ReadPerson', personId);
      result = JSON.parse(result.toString());
      res.json(result);
    }
    catch (err) {
        res.json(err.message)
    } finally {
        closeConnetion(gateway);
    }
});




module.exports = router;