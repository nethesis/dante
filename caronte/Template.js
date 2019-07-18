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
const mjml2html = require("mjml");

class Template {
  constructor() {
    this.htmlOutput = null;
  }

  async create(url, parsedQuery) {
    // extract query
    var theme = parsedQuery["?theme"];
    var last = parsedQuery.last;
    var lang = parsedQuery.lang;

    // configure lang
    i18n.configure({
      locales: [lang],
      directory: process.pkg
        ? "/usr/share/dante/beatrice/i18n/" + lang + "/"
        : __dirname + "/../beatrice/public/i18n/" + lang + "/",
      objectNotation: true
    });
    moment.locale(lang);

    // create machine string
    var machineString =
      i18n.__("caronte.machine") +
      ": <code><b>" +
      os.hostname() +
      "</b></code>";

    // create translation string
    var htmlString = "";
    switch (last) {
      case "week":
        htmlString +=
          i18n.__("caronte.last_week") +
          ": <b>" +
          moment()
            .subtract(7, "days")
            .format("DD MMM YYYY") +
          "</b> - <b>" +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY") +
          "</b>";
        break;

      case "month":
        htmlString +=
          i18n.__("caronte.last_month") +
          ": <b>" +
          moment()
            .subtract(1, "months")
            .format("DD MMM YYYY") +
          "</b> - <b>" +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY") +
          "</b>";
        break;

      case "halfyear":
        htmlString +=
          i18n.__("caronte.last_halfyear") +
          ": <b>" +
          moment()
            .subtract(6, "months")
            .format("DD MMM YYYY") +
          "</b> - <b>" +
          moment()
            .subtract(1, "days")
            .format("DD MMM YYYY") +
          "</b>";
        break;
    }

    this.htmlOutput = mjml2html(
      `
      <mjml>

      <mj-head>
        <mj-attributes>
            <mj-all font-family="Arial" />
            <mj-class name="mjclass" color="` +
        (theme == "dark" ? "#1d1e1e" : "#ffffff") +
        `"/>
        </mj-attributes>

        <mj-style>
          .background-wrapper {
            background: linear-gradient(0deg, #2e2e2e 60%, #e0e1e2 50%);
          }

          .shadow-down {
            box-shadow: 0px 0px 20px 0px rgba(0, 0, 0, 0.15);
          }
        </mj-style>
      </mj-head>

      <mj-body css-class="background-wrapper">
        <mj-section css-class="shadow-down" background-color="` +
        (theme == "dark" ? "#ffffff" : "#1d1e1e") +
        `">
          <mj-column>
            <mj-text mj-class="mjclass" align="center" font-weight="bold" font-size="25px">` +
        i18n.__("caronte.report_title") +
        ` 🎉</mj-text>

          <mj-text mj-class="mjclass" align="center" font-size="80px">📊</mj-text>

            <mj-text mj-class="mjclass" align="center" font-size="16px">` +
        htmlString +
        `</mj-text>
            <mj-text mj-class="mjclass" align="center" font-size="16px">
                ` +
        machineString +
        `
            </mj-text>
            <mj-spacer height="10px" />
            <mj-divider padding="0px 0px" border-width="1px" border-style="solid" border-color="` +
        (theme == "light" ? "#c0c1c3" : "#2e2e2e") +
        `" />

            <mj-divider padding="0px 0px" border-width="1px" border-style="solid" border-color="` +
        (theme == "light" ? "#c0c1c3" : "#2e2e2e") +
        `" />
            <mj-spacer height="10px" />
            <mj-button background-color="#f45e43" color="white" href="` +
        url +
        `">
              ` +
        i18n.__("caronte.view_more") +
        `
            </mj-button>
          </mj-column>
        </mj-section>

      </mj-body>
    </mjml>
  `,
      {}
    );

    return this.htmlOutput.html;
  }
}
module.exports = Template;
