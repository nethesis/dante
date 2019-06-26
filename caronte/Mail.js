const os = require("os");
const i18n = require("i18n");
const moment = require("moment");
const nodemailer = require("nodemailer");

class Mail {
  constructor() {}

  async send(html, addresses, url, parsedQuery, pdf) {
    // extract query
    var last = parsedQuery.last;
    var lang = parsedQuery.lang;

    // initialize transporter
    var transporter = nodemailer.createTransport({
      host: "localhost",
      port: 587,
      secure: false,
      ignoreTLS: true
    });

    // configure lang
    i18n.configure({
      locales: [lang],
      directory: process.pkg
        ? "/usr/share/dante/beatrice/i18n/"
        : __dirname + "/../beatrice/public/i18n/",
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
      from: '"Root" <root@' + os.hostname() + ">",
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
      html: html,
      attachments: [
        {
          filename: last + "ly_report.pdf",
          content: pdf.toString("base64"),
          encoding: "base64"
        }
      ]
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
