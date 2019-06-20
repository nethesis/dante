const nodemailer = require("nodemailer");

class Mail {
  constructor() {}

  async send(html, addresses) {
    // initialize transporter
    var transporter = nodemailer.createTransport({
      host: "localhost",
      port: 587,
      secure: false,
      ignoreTLS: true
    });

    // setup e-mail data with unicode symbols
    var mailOptions = {
      from: '"Fred Foo ?" <foo@blurdybloop.com>',
      to: addresses,
      subject: "📈 Weekly Report | 09 Jun 2019 - 15 Jun 2019",
      text: "Hello world ?",
      html: html
    };

    // send mail with defined transport object
    return await transporter.sendMail(mailOptions);
  }
}
module.exports = Mail;
