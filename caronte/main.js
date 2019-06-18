const puppeteer = require("puppeteer");

let _browser;
let _page;

let args = process.argv.slice(2);

puppeteer
  .launch({
    args: ["--no-sandbox"],
    defaultViewport: {
      width: 1024,
      height: 768
    },
    executablePath: process.pkg
      ? puppeteer.executablePath().replace(__dirname, ".")
      : null
  })
  .then(browser => (_browser = browser))
  .then(browser => (_page = browser.newPage()))
  .then(page => page.goto(args[0])) // args[0] contains beatrice url
  .then(() => _page)
  .then(page => page.waitFor(2000)) // wait for charts animations
  .then(() => _page)
  .then(page =>
    page.screenshot({
      path: "caronte-" + new Date().toISOString() + ".png",
      fullPage: true
    })
  )
  .then(() => _browser.close());
