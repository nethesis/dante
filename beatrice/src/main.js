import Vue from "vue";

import VueResource from "vue-resource";
Vue.use(VueResource);

import VueI18n from "vue-i18n";
Vue.use(VueI18n);

import App from "./App.vue";
import router from "./router";

window.$ = window.jQuery = require("jquery");

require("semantic-ui-css/semantic.min.css");
require("semantic-ui-css/semantic.min.js");

import VueApexCharts from "vue-apexcharts";
Vue.component("apexchart", VueApexCharts);

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
  render: h => h(App)
});

app.$http.get("./i18n/" + lang + ".json").then(
  success => {
    i18n.setLocaleMessage(lang, success.body);
    app.$mount("#app");
  },
  error => {
    console.error(error);
  }
);
