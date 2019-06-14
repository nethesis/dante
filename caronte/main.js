const fs = require("fs");
const puppeteer = require("puppeteer");

let _browser;
let _page;

const config = JSON.parse(fs.readFileSync("./config.json", "utf-8"));

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
  .then(page => page.goto(config.beatrice_url))
  .then(() => _page)
  .then(page => {
    if (config.theme == "dark") {
      page.click("#toggleTheme");
    }
  })
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

// yum install -y at-spi2-atk gtk3 libXScrnSaver
