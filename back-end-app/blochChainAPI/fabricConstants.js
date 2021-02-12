var path = require('path');

const channelName = 'mychannel';
const chaincodeName = 'basic';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser9';

exports.channelName = channelName;
exports.chaincodeName = chaincodeName;
exports.mspOrg1 = mspOrg1;
exports.walletPath = walletPath;
exports.org1UserId = org1UserId;