#!/usr/bin/env node
const urlparse = require("url-parse");
const querystring = require("querystring");

const Screenshot = require("./Screenshot");
const Template = require("./Template");
const Mail = require("./Mail");

(async () => {
  // read args
  var args = process.argv.slice(2);

  // check args
  if (!args[0] || args[0].length == 0) {
    console.error("Beatrice 'URL' is required.");
    process.exit(1);
  }

  if (!args[1] || args[1].length == 0) {
    console.error("Sender address[es]' are required.");
    process.exit(1);
  }

  // get args
  var url = args[0];
  var addresses = args[1];

  // parse url
  var queryStringParsed = new urlparse(url.replace("/#/?", "?")).query;
  var parsedQuery = querystring.parse(queryStringParsed);

  // process images and send mail
  const screenshot = new Screenshot();
  var image = await screenshot.shot(url);

  const template = new Template();
  var html = await template.create(image, url, parsedQuery);

  const mail = new Mail();
  var status = await mail.send(html, addresses, url, parsedQuery);

  process.exit(0);
})();
