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
const os = require("os");
const i18n = require("i18n");
const moment = require("moment");
const nodemailer = require("nodemailer");

class Mail {
  constructor() {}

  async send(html, addresses, url, parsedQuery, from) {
    // extract query
    var last = parsedQuery.last;
    var lang = parsedQuery.lang;

    // initialize transporter
    var transporter = nodemailer.createTransport({
      host: "127.0.0.1",
      port: 25,
      secure: false,
      ignoreTLS: true
    });

    // configure lang
    i18n.configure({
      locales: [lang],
      directory: process.pkg
        ? "/usr/share/dante/beatrice/i18n/" + lang + "/"
        : __dirname + "/../beatrice/public/i18n/" + lang + "/",
      objectNotation: true
    });
    moment.locale(lang);

    var dateRange = "";
    switch (last) {
      case "week":
        dateRange =
          moment()
            .subtract(7, "days")
            .format("DD MMM YYYY") +
          " - " +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY");
        break;

      case "month":
        dateRange =
          moment()
            .subtract(1, "months")
            .format("DD MMM YYYY") +
          " - " +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY");
        break;

      case "halfyear":
        dateRange =
          moment()
            .subtract(6, "months")
            .format("DD MMM YYYY") +
          " - " +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY");
        break;
    }

    // setup e-mail data with unicode symbols
    var mailOptions = {
      from: from,
      to: addresses,
      subject:
        "📈 " +
        i18n.__("caronte." + last + "ly") +
        " (" +
        os.hostname() +
        ") | " +
        dateRange,
      text:
        i18n.__("caronte." + last + "ly") +
        " (" +
        os.hostname() +
        ") link: " +
        url,
      html: html
    };

    // send mail with defined transport object
    try {
      return await transporter.sendMail(mailOptions);
    } catch (error) {
      console.error(error);
      process.exit(1);
    }
  }
}
module.exports = Mail;
