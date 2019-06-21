/*
 * Copyright (C) 2019 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dante project.
 *
 * Dante is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dante is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dante.  If not, see COPYING.
 */
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
    try {
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
      await this.page.goto(url, { waitUntil: "networkidle2", timeout: 60000 });

      // wait charts animations
      await this.page.waitFor(2000);

      // take screenshot
      var image = await this.screenshotDOMElement({
        selector: ".ui.segment",
        padding: -1,
        height: 768
      });

      // take pdf
      await this.page.emulateMedia("screen");
      var pdf = await this.page.pdf({
        format: "Tabloid",
        printBackground: true
      });

      // return image and pdf
      return {
        image: image,
        pdf: pdf
      };
    } catch (error) {
      // close browser instance and exit process
      console.error(error);
      await this.browser.close();
      process.exit(1);
    }
  }
}
module.exports = Screenshot;
