const puppeteer = require("puppeteer");

class Screenshot {
  constructor() {
    this.browser = null;
    this.page = null;
  }

  async screenshotDOMElement(opts = {}) {
    const padding = "padding" in opts ? opts.padding : 0;
    const path = "path" in opts ? opts.path : null;
    const selector = opts.selector;

    if (!selector) throw Error("Please provide a selector.");

    const rect = await this.page.evaluate(selector => {
      const element = document.querySelector(selector);
      if (!element) return null;
      const { x, y, width, height } = element.getBoundingClientRect();
      return { left: x, top: y, width, height, id: element.id };
    }, selector);

    if (!rect)
      throw Error(`Could not find element that matches selector: ${selector}.`);

    return await this.page.screenshot({
      clip: {
        x: rect.left - padding,
        y: rect.top - padding,
        width: rect.width + padding * 2,
        height: "height" in opts ? opts.height : rect.height + padding * 2
      },
      encoding: "base64"
    });
  }

  async shot(url) {
    // create browser instance
    this.browser = await puppeteer.launch({
      defaultViewport: {
        width: 1024,
        height: 0
      },
      executablePath: process.pkg
        ? puppeteer.executablePath().replace(__dirname, ".")
        : null,
      args: ["--no-sandbox"]
    });

    // create page
    this.page = await this.browser.newPage();

    // open url
    await this.page.goto(url, { waitUntil: "networkidle2" });

    // wait charts animations
    await this.page.waitFor(2000);

    // take screenshot
    return await this.screenshotDOMElement({
      selector: ".ui.segment",
      padding: -1,
      height: 768
    });
  }
}
module.exports = Screenshot;
