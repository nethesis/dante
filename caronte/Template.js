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

  async create(image, url, parsedQuery) {
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
            <mj-image width="100px" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAYAAABccqhmAAAACXBIWXMAAAsTAAALEwEAmpwYAABCQWlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4KPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iQWRvYmUgWE1QIENvcmUgNS42LWMwMTQgNzkuMTU2Nzk3LCAyMDE0LzA4LzIwLTA5OjUzOjAyICAgICAgICAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iCiAgICAgICAgICAgIHhtbG5zOmRjPSJodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyIKICAgICAgICAgICAgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iCiAgICAgICAgICAgIHhtbG5zOnN0RXZ0PSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VFdmVudCMiCiAgICAgICAgICAgIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIgogICAgICAgICAgICB4bWxuczpwaG90b3Nob3A9Imh0dHA6Ly9ucy5hZG9iZS5jb20vcGhvdG9zaG9wLzEuMC8iCiAgICAgICAgICAgIHhtbG5zOnRpZmY9Imh0dHA6Ly9ucy5hZG9iZS5jb20vdGlmZi8xLjAvIgogICAgICAgICAgICB4bWxuczpleGlmPSJodHRwOi8vbnMuYWRvYmUuY29tL2V4aWYvMS4wLyI+CiAgICAgICAgIDx4bXA6Q3JlYXRlRGF0ZT4yMDE2LTAyLTI3VDA4OjQ1OjM5WjwveG1wOkNyZWF0ZURhdGU+CiAgICAgICAgIDx4bXA6TW9kaWZ5RGF0ZT4yMDE2LTAyLTI3VDEzOjMxOjIwKzA0OjAwPC94bXA6TW9kaWZ5RGF0ZT4KICAgICAgICAgPHhtcDpNZXRhZGF0YURhdGU+MjAxNi0wMi0yN1QxMzozMToyMCswNDowMDwveG1wOk1ldGFkYXRhRGF0ZT4KICAgICAgICAgPHhtcDpDcmVhdG9yVG9vbD5BZG9iZSBQaG90b3Nob3AgQ0MgMjAxNCAoTWFjaW50b3NoKTwveG1wOkNyZWF0b3JUb29sPgogICAgICAgICA8ZGM6Zm9ybWF0PmltYWdlL3BuZzwvZGM6Zm9ybWF0PgogICAgICAgICA8eG1wTU06SGlzdG9yeT4KICAgICAgICAgICAgPHJkZjpTZXE+CiAgICAgICAgICAgICAgIDxyZGY6bGkgcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6YWN0aW9uPmNvbnZlcnRlZDwvc3RFdnQ6YWN0aW9uPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6cGFyYW1ldGVycz5mcm9tIGltYWdlL3BuZyB0byBhcHBsaWNhdGlvbi92bmQuYWRvYmUucGhvdG9zaG9wPC9zdEV2dDpwYXJhbWV0ZXJzPgogICAgICAgICAgICAgICA8L3JkZjpsaT4KICAgICAgICAgICAgICAgPHJkZjpsaSByZGY6cGFyc2VUeXBlPSJSZXNvdXJjZSI+CiAgICAgICAgICAgICAgICAgIDxzdEV2dDphY3Rpb24+c2F2ZWQ8L3N0RXZ0OmFjdGlvbj4KICAgICAgICAgICAgICAgICAgPHN0RXZ0Omluc3RhbmNlSUQ+eG1wLmlpZDo0NDNFRTE5RTBCMjA2ODExODIyQUM1MTg3MTkxRkQ2Qzwvc3RFdnQ6aW5zdGFuY2VJRD4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OndoZW4+MjAxNi0wMi0yN1QxMzoxOToyMyswNDowMDwvc3RFdnQ6d2hlbj4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OmNoYW5nZWQ+Lzwvc3RFdnQ6Y2hhbmdlZD4KICAgICAgICAgICAgICAgPC9yZGY6bGk+CiAgICAgICAgICAgICAgIDxyZGY6bGkgcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6YWN0aW9uPmNvbnZlcnRlZDwvc3RFdnQ6YWN0aW9uPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6cGFyYW1ldGVycz5mcm9tIGltYWdlL3BuZyB0byBhcHBsaWNhdGlvbi92bmQuYWRvYmUucGhvdG9zaG9wPC9zdEV2dDpwYXJhbWV0ZXJzPgogICAgICAgICAgICAgICA8L3JkZjpsaT4KICAgICAgICAgICAgICAgPHJkZjpsaSByZGY6cGFyc2VUeXBlPSJSZXNvdXJjZSI+CiAgICAgICAgICAgICAgICAgIDxzdEV2dDphY3Rpb24+c2F2ZWQ8L3N0RXZ0OmFjdGlvbj4KICAgICAgICAgICAgICAgICAgPHN0RXZ0Omluc3RhbmNlSUQ+eG1wLmlpZDo0NTNFRTE5RTBCMjA2ODExODIyQUM1MTg3MTkxRkQ2Qzwvc3RFdnQ6aW5zdGFuY2VJRD4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OndoZW4+MjAxNi0wMi0yN1QxMzoxOToyMyswNDowMDwvc3RFdnQ6d2hlbj4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OmNoYW5nZWQ+Lzwvc3RFdnQ6Y2hhbmdlZD4KICAgICAgICAgICAgICAgPC9yZGY6bGk+CiAgICAgICAgICAgICAgIDxyZGY6bGkgcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6YWN0aW9uPnNhdmVkPC9zdEV2dDphY3Rpb24+CiAgICAgICAgICAgICAgICAgIDxzdEV2dDppbnN0YW5jZUlEPnhtcC5paWQ6ZDljYzBlY2MtNDYxMi00NmU1LWE2MzAtZDUxNWJhMGNhZTA2PC9zdEV2dDppbnN0YW5jZUlEPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6d2hlbj4yMDE2LTAyLTI3VDEzOjMxOjIwKzA0OjAwPC9zdEV2dDp3aGVuPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6c29mdHdhcmVBZ2VudD5BZG9iZSBQaG90b3Nob3AgQ0MgMjAxNCAoTWFjaW50b3NoKTwvc3RFdnQ6c29mdHdhcmVBZ2VudD4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OmNoYW5nZWQ+Lzwvc3RFdnQ6Y2hhbmdlZD4KICAgICAgICAgICAgICAgPC9yZGY6bGk+CiAgICAgICAgICAgICAgIDxyZGY6bGkgcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6YWN0aW9uPmNvbnZlcnRlZDwvc3RFdnQ6YWN0aW9uPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6cGFyYW1ldGVycz5mcm9tIGFwcGxpY2F0aW9uL3ZuZC5hZG9iZS5waG90b3Nob3AgdG8gaW1hZ2UvcG5nPC9zdEV2dDpwYXJhbWV0ZXJzPgogICAgICAgICAgICAgICA8L3JkZjpsaT4KICAgICAgICAgICAgICAgPHJkZjpsaSByZGY6cGFyc2VUeXBlPSJSZXNvdXJjZSI+CiAgICAgICAgICAgICAgICAgIDxzdEV2dDphY3Rpb24+ZGVyaXZlZDwvc3RFdnQ6YWN0aW9uPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6cGFyYW1ldGVycz5jb252ZXJ0ZWQgZnJvbSBhcHBsaWNhdGlvbi92bmQuYWRvYmUucGhvdG9zaG9wIHRvIGltYWdlL3BuZzwvc3RFdnQ6cGFyYW1ldGVycz4KICAgICAgICAgICAgICAgPC9yZGY6bGk+CiAgICAgICAgICAgICAgIDxyZGY6bGkgcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6YWN0aW9uPnNhdmVkPC9zdEV2dDphY3Rpb24+CiAgICAgICAgICAgICAgICAgIDxzdEV2dDppbnN0YW5jZUlEPnhtcC5paWQ6NmY2OTkzN2YtYTI0NS00OTI2LWIyZGUtODdhMDk5YzAwMTkwPC9zdEV2dDppbnN0YW5jZUlEPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6d2hlbj4yMDE2LTAyLTI3VDEzOjMxOjIwKzA0OjAwPC9zdEV2dDp3aGVuPgogICAgICAgICAgICAgICAgICA8c3RFdnQ6c29mdHdhcmVBZ2VudD5BZG9iZSBQaG90b3Nob3AgQ0MgMjAxNCAoTWFjaW50b3NoKTwvc3RFdnQ6c29mdHdhcmVBZ2VudD4KICAgICAgICAgICAgICAgICAgPHN0RXZ0OmNoYW5nZWQ+Lzwvc3RFdnQ6Y2hhbmdlZD4KICAgICAgICAgICAgICAgPC9yZGY6bGk+CiAgICAgICAgICAgIDwvcmRmOlNlcT4KICAgICAgICAgPC94bXBNTTpIaXN0b3J5PgogICAgICAgICA8eG1wTU06RGVyaXZlZEZyb20gcmRmOnBhcnNlVHlwZT0iUmVzb3VyY2UiPgogICAgICAgICAgICA8c3RSZWY6aW5zdGFuY2VJRD54bXAuaWlkOmQ5Y2MwZWNjLTQ2MTItNDZlNS1hNjMwLWQ1MTViYTBjYWUwNjwvc3RSZWY6aW5zdGFuY2VJRD4KICAgICAgICAgICAgPHN0UmVmOmRvY3VtZW50SUQ+eG1wLmRpZDo0NDNFRTE5RTBCMjA2ODExODIyQUM1MTg3MTkxRkQ2Qzwvc3RSZWY6ZG9jdW1lbnRJRD4KICAgICAgICAgICAgPHN0UmVmOm9yaWdpbmFsRG9jdW1lbnRJRD54bXAuZGlkOjQ0M0VFMTlFMEIyMDY4MTE4MjJBQzUxODcxOTFGRDZDPC9zdFJlZjpvcmlnaW5hbERvY3VtZW50SUQ+CiAgICAgICAgIDwveG1wTU06RGVyaXZlZEZyb20+CiAgICAgICAgIDx4bXBNTTpEb2N1bWVudElEPmFkb2JlOmRvY2lkOnBob3Rvc2hvcDo0OTBkODgzOC0xZGMzLTExNzktOWQzYS1jZWI1YzU4YzRiZDQ8L3htcE1NOkRvY3VtZW50SUQ+CiAgICAgICAgIDx4bXBNTTpJbnN0YW5jZUlEPnhtcC5paWQ6NmY2OTkzN2YtYTI0NS00OTI2LWIyZGUtODdhMDk5YzAwMTkwPC94bXBNTTpJbnN0YW5jZUlEPgogICAgICAgICA8eG1wTU06T3JpZ2luYWxEb2N1bWVudElEPnhtcC5kaWQ6NDQzRUUxOUUwQjIwNjgxMTgyMkFDNTE4NzE5MUZENkM8L3htcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD4KICAgICAgICAgPHBob3Rvc2hvcDpDb2xvck1vZGU+MzwvcGhvdG9zaG9wOkNvbG9yTW9kZT4KICAgICAgICAgPHRpZmY6T3JpZW50YXRpb24+MTwvdGlmZjpPcmllbnRhdGlvbj4KICAgICAgICAgPHRpZmY6WFJlc29sdXRpb24+NzIwMDAwLzEwMDAwPC90aWZmOlhSZXNvbHV0aW9uPgogICAgICAgICA8dGlmZjpZUmVzb2x1dGlvbj43MjAwMDAvMTAwMDA8L3RpZmY6WVJlc29sdXRpb24+CiAgICAgICAgIDx0aWZmOlJlc29sdXRpb25Vbml0PjI8L3RpZmY6UmVzb2x1dGlvblVuaXQ+CiAgICAgICAgIDxleGlmOkNvbG9yU3BhY2U+NjU1MzU8L2V4aWY6Q29sb3JTcGFjZT4KICAgICAgICAgPGV4aWY6UGl4ZWxYRGltZW5zaW9uPjI1NjwvZXhpZjpQaXhlbFhEaW1lbnNpb24+CiAgICAgICAgIDxleGlmOlBpeGVsWURpbWVuc2lvbj4yNTY8L2V4aWY6UGl4ZWxZRGltZW5zaW9uPgogICAgICA8L3JkZjpEZXNjcmlwdGlvbj4KICAgPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAKPD94cGFja2V0IGVuZD0idyI/PlV2eUsAAAAgY0hSTQAAeiUAAICDAAD5/wAAgOkAAHUwAADqYAAAOpgAABdvkl/FRgAAHHFJREFUeNrsnXl0XdV97z930mxLtjwx2Rhj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" />
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

            <mj-image fluid-on-mobile="true" padding="0px 0px" src="data:image/png;base64,` +
        image +
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
