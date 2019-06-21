const os = require("os");
const moment = require("moment");
const nodemailer = require("nodemailer");

class Mail {
  constructor() {}

  async send(html, addresses, url, parsedQuery, pdf) {
    // extract query
    var last = parsedQuery.last;

    // initialize transporter
    var transporter = nodemailer.createTransport({
      host: "localhost",
      port: 587,
      secure: false,
      ignoreTLS: true
    });

    var dateRange = "";
    switch (last) {
      case "day":
        dateRange = moment()
          .subtract(1, "days")
          .format("LL");
        break;

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
    }

    // setup e-mail data with unicode symbols
    var mailOptions = {
      from: '"Root" <root@' + os.hostname() + ">",
      to: addresses,
      subject:
        "📈 " +
        last.charAt(0).toUpperCase() +
        last.slice(1) +
        "ly Report (" +
        os.hostname() +
        ") | " +
        dateRange,
      text:
        last.charAt(0).toUpperCase() +
        last.slice(1) +
        "ly Report (" +
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
    return await transporter.sendMail(mailOptions);
  }
}
module.exports = Mail;
