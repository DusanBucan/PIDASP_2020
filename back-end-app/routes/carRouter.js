var express = require('express');
var router = express.Router();

var app = require("../app");
var {getConection, closeConnetion, prettyJSONString} = require("../blochChainAPI/blockChainConnection.js");
const { chaincodeName, channelName, org1UserId } = require('../blochChainAPI/fabricConstants.js');


router.get('/all', async (req, res, next) => {

    let wallet = req.app.get("fabricWallet");
    let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
    try {
      let result = await contract.evaluateTransaction('GetAllCars');
      result = JSON.parse(result.toString());
      res.json(result); 
    }
    catch (err) {
        res.json(err.message)
    } finally {
        closeConnetion(gateway);
    }
});

router.get('/:carId', async (req, res, next) => {
    const carId = req.params.carId;
    let wallet = req.app.get("fabricWallet");
    let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
    
    try {
      result = await contract.evaluateTransaction('ReadCar', carId);
      result = JSON.parse(result.toString());
      res.json(result);
    }
    catch (err) {
        res.json(err.message)
    } finally {
        closeConnetion(gateway);
    }
});


router.get('/filterColor/:color',  async (req, res, next) => {

  const color = req.params.color;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  try {
    result = await contract.evaluateTransaction('GetAllCarsByCollor', color);
    result = JSON.parse(result.toString());
    res.json(result);
  } catch (err) {
    res.json(err.message)
  }
  finally {
    closeConnetion(gateway);
  }
} )

router.get('/filterColor/:color/:owner',  async (req, res, next) => {

  const color = req.params.color;
  const owner = req.params.owner;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  try {
    result = await contract.evaluateTransaction('GetAllCarsByCollorAndOwner', color, owner);
    result = JSON.parse(result.toString());
    res.json(result);
  } catch (err) {
    res.json(err.message)
  }
  finally {
    closeConnetion(gateway);
  }
} )


router.get('/erros/:carId',  async (req, res, next) => {

    const carId = req.params.carId;
    let wallet = req.app.get("fabricWallet");
    let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
    try {
      result = await contract.evaluateTransaction('ReadAllCarBreakDown', carId);
      result = JSON.parse(result.toString());
      res.json(result);
    } catch (err) {
      res.json(err.message)
    }
    finally {
      closeConnetion(gateway);
    }
} )





router.post('/makeBreakdown', async (req, res, next)=> {
  const createCarBrakedownData = req.body;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  try {
    result = await contract.submitTransaction('CreateCarBrakedown', 
    createCarBrakedownData.description, 
    createCarBrakedownData.price, 
    createCarBrakedownData.carId
    );
      console.log('*** Result: committed');
    if (`${result}` !== '') {
      result = JSON.parse(result.toString());
    } else {
      result = "";
    }
    res.json(result);
  } catch (err) {
      res.json(err);
  } finally {
      closeConnetion(gateway);
  }

});

router.post('/fixBreakdown', async (req, res, next)=> {
  const fixCarBrakedownData = req.body;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  try {
    result = await contract.submitTransaction('FixCarBrakedown', 
      fixCarBrakedownData.id, 
      fixCarBrakedownData.mechanicId
    );
    console.log('*** Result: committed');
    if (`${result}` !== '') {
      result = JSON.parse(result.toString());
    } else {
      result = "";
    }
    res.json(result);
  } catch (err) {
    res.json(err.message)
  } finally {
    closeConnetion(gateway);
  }
});


router.post('/changeOwner', async (req, res, next)=> {
  const changeCarOwnerData = req.body;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  try {
      result = await contract.submitTransaction('UpdateCarOwner', 
      changeCarOwnerData.ID, 
      changeCarOwnerData.newOwnerId, 
      changeCarOwnerData.buyWithErrors
    );
    console.log('*** Result: committed');
    if (`${result}` !== '') {
      result = JSON.parse(result.toString());
    } else {
      result = "";
    }
    res.json(result);
  } catch (err) {
    res.json(err.message)
  }
  finally {
    closeConnetion(gateway);
  }
 
})

router.post('/changeColor', async (req, res, next)=> {
  const changeColorData = req.body;
  let wallet = req.app.get("fabricWallet");
  let {contract, gateway} = await getConection(wallet,channelName, chaincodeName, org1UserId);
  
  try {
    result = await contract.submitTransaction('UpdateCarColor', 
    changeColorData.ID, changeColorData.Color,
    changeColorData.Cost, changeColorData.mechanicId)
    console.log('*** Result: committed');
    if (`${result}` !== '') {
      result = JSON.parse(result.toString());
    } else {
      result = "";
      res.json(result);
    }
  } catch (err) {
       res.json(err.message)
  } finally {
      closeConnetion(gateway);
  }
})

module.exports = router;