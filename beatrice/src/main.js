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
import Vue from "vue";

import "./filters";

import VueResource from "vue-resource";
Vue.use(VueResource);

import VueI18n from "vue-i18n";
Vue.use(VueI18n);

import App from "./App.vue";
import router from "./router";

import VueApexCharts from "vue-apexcharts";
Vue.component("apexchart", VueApexCharts);

window.$ = window.jQuery = require("jquery");

require("semantic-ui-css/semantic.min.css");
require("semantic-ui-css/semantic.min.js");

Vue.config.productionTip = false;

var url = new URL(window.location.href.replace("/#/?", "?"));
var params = new URLSearchParams(url.search);

const lang = params.get("lang") || "en";
const i18n = new VueI18n({
  locale: lang
});

var app = new Vue({
  i18n,
  router,
  apiHost: "your_virgilio_ip",
  render: h => h(App)
});

app.$http.get(app.$options.apiHost + "/lang/" + lang).then(
  success => {
    i18n.setLocaleMessage(lang, success.body.lang);
    app.$mount("#app");
  },
  error => {
    console.error(error);
  }
);