<!--
Copyright (C) 2019 Nethesis S.r.l.
http://www.nethesis.it - info@nethesis.it
 This file is part of Dante project.

 Dante is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published by
 the Free Software Foundation, either version 3 of the License,
 or any later version.

 Dante is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with Dante.  If not, see COPYING.
-->
<template>
  <div id="app">
    <div v-if="!isMobile" :class="['ui pointing menu', lightTheme ? '' : 'inverted']">
      <a href="#/" :class="[getCurrentPath('') ? 'active' : '', 'item']">{{$t('home.title')}}</a>
      <div class="item">
        <div class="ui buttons" :class="lightTheme ? '' : 'inverted'">
          <span
            :data-tooltip="maxDays < 7 ? $t('home.data_not_available_in_range') : null"
            data-position="bottom center"
          >
            <button
              :disabled="maxDays < 7"
              @click="setFilterDate('week')"
              class="ui button"
              :class="[lightTheme ? '' : 'inverted', filterDate == 'week' && !showCustomInterval ? 'active' : '']"
            >{{$t('home.last_week')}}</button>
          </span>
          <span
            :data-tooltip="maxDays < 31 ? $t('home.data_not_available_in_range') : null"
            data-position="bottom center"
          >
            <button
              :disabled="maxDays < 31"
              @click="setFilterDate('month')"
              class="ui button"
              :class="[lightTheme ? '' : 'inverted', filterDate == 'month' && !showCustomInterval ? 'active' : '']"
            >{{$t('home.last_month')}}</button>
          </span>
          <span
            :data-tooltip="maxDays < 181 ? $t('home.data_not_available_in_range') : null"
            data-position="bottom center"
          >
            <button
              :disabled="maxDays < 181"
              @click="setFilterDate('halfyear')"
              class="ui button"
              :class="[lightTheme ? '' : 'inverted', filterDate == 'halfyear' && !showCustomInterval ? 'active' : '']"
            >{{$t('home.last_halfyear')}}</button>
          </span>
          <span data-position="bottom center">
            <button
              @click="setCustomInterval()"
              class="ui button mg-left-15"
              :class="[lightTheme ? '' : 'inverted', showCustomInterval ? 'active' : '']"
            >{{$t('home.custom_interval')}}</button>
          </span>
        </div>
        <div v-show="showCustomInterval" class="customIntervalPanel">
          <datepicker
            :class="['datepicker', lightTheme ? '' : 'inverted', customIntervalError ? 'customIntervalError' : '']"
            :placeholder="$t('home.start_date')"
            v-model="customStartDate"
            :disabled-dates="disabledDates"
            id="customStartDate"
          ></datepicker>
          <span class="mg-left-5">{{$t('home.to')}}</span>
          <datepicker
            :class="['datepicker', lightTheme ? '' : 'inverted', customIntervalError ? 'customIntervalError' : '']"
            :placeholder="$t('home.end_date')"
            v-model="customEndDate"
            :disabled-dates="disabledDates"
            id="customEndDate"
          ></datepicker>
        </div>
      </div>
      <div class="right menu">
        <div class="item">
          <div class="ui transparent icon input" :class="lightTheme ? '' : 'inverted'">
            <input v-model="searchString" type="text" :placeholder="$t('home.search')+'...'" />
            <i
              @click="searchString.length > 0 ? searchString = '' : undefined"
              :class="[searchString.length > 0 ? 'remove link' : 'search', 'icon']"
            ></i>
          </div>
        </div>
        <div class="item">
          <div class="ui buttons">
            <div
              id="toggleTheme"
              @click="setTheme()"
              class="ui button"
              :class="lightTheme ? 'black' : 'inverted'"
            >{{lightTheme ? $t('home.dark_theme') : $t('home.light_theme')}}</div>
          </div>
          <div class="ui compact menu mg-left-10" :class="lightTheme ? '' : 'inverted'">
            <div class="ui simple dropdown item">
              {{$t('home.colors')}}
              <i class="dropdown icon"></i>
              <div class="menu">
                <div
                  @click="setPalette(1)"
                  class="item"
                  :class="colorPalette == 'palette1' ? 'selected' : ''"
                >{{$t('home.palette')}} 1</div>
                <div
                  @click="setPalette(2)"
                  class="item"
                  :class="colorPalette == 'palette2' ? 'selected' : ''"
                >{{$t('home.palette')}} 2</div>
                <div
                  @click="setPalette(3)"
                  class="item"
                  :class="colorPalette == 'palette3' ? 'selected' : ''"
                >{{$t('home.palette')}} 3</div>
                <div
                  @click="setPalette(4)"
                  class="item"
                  :class="colorPalette == 'palette4' ? 'selected' : ''"
                >{{$t('home.palette')}} 4</div>
                <div
                  @click="setPalette(5)"
                  class="item"
                  :class="colorPalette == 'palette5' ? 'selected' : ''"
                >{{$t('home.palette')}} 5</div>
                <div
                  @click="setPalette(6)"
                  class="item"
                  :class="colorPalette == 'palette6' ? 'selected' : ''"
                >{{$t('home.palette')}} 6</div>
                <div
                  @click="setPalette(7)"
                  class="item"
                  :class="colorPalette == 'palette7' ? 'selected' : ''"
                >{{$t('home.palette')}} 7</div>
                <div
                  @click="setPalette(8)"
                  class="item"
                  :class="colorPalette == 'palette8' ? 'selected' : ''"
                >{{$t('home.palette')}} 8</div>
                <div
                  @click="setPalette(9)"
                  class="item"
                  :class="colorPalette == 'palette9' ? 'selected' : ''"
                >{{$t('home.palette')}} 9</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div :class="['ui segment', lightTheme ? '' : 'inverted']">
      <router-view />
    </div>
  </div>
</template>

<script>
import Datepicker from 'vuejs-datepicker';
var moment = require("moment");

export default {
  name: "home",
  mounted() {
    document.title =
      this.$t("home.title") +
      ": " +
      this.$t("caronte.last_" + (this.$route.query.last || "week"));
  },
  data() {
    // set locale
    moment.locale(this.$options.lang);

    return {
      lightTheme: this.$route.query.theme
        ? this.$route.query.theme == "light"
        : true,
      colorPalette: this.$route.query.palette || "palette1",
      filterDate: this.$route.query.last || "week",
      language: this.$route.query.lang || "en",
      searchString: "",
      isMobile: /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
        navigator.userAgent
      ),
      maxDays: 0,
      customStartDate: this.getDate(this.$route.query.customStartDate),
      customEndDate: this.initCustomEndDate(),
      disabledDates: {
        from: moment().toDate(),
        to: moment().startOf("day").subtract(181, "days").toDate()
      },
      showCustomInterval: this.$route.query.last === 'custom',
      customIntervalError: false
    };
  },
  components: {
  	Datepicker
  },
  watch: {
    customStartDate: function() {
      this.setCustomInterval();
    },
    customEndDate: function() {
      this.setCustomInterval();
    }
  },
  methods: {
    getDate(dateString) {
      var momentDate = moment(dateString, 'YYYY-MM-DD');

      if (momentDate && momentDate.toDate() != 'Invalid Date') {
        return momentDate.toDate();
      } else {
        return null;
      }
    },
    initCustomEndDate() {
      var date = this.getDate(this.$route.query.customEndDate);

      if (date === null) {
        // return today as default value
        return moment().startOf("day").toDate();
      } else {
        return date;
      }
    },
    setTheme() {
      this.lightTheme = !this.lightTheme;
      this.updateQuery();
    },
    setPalette(paletteNumber) {
      this.colorPalette = "palette" + paletteNumber;
      this.updateQuery();
    },
    setFilterDate(last) {
      this.filterDate = last;
      this.showCustomInterval = false;
      document.title =
        this.$i18n.t("home.title") +
        ": " +
        this.$i18n.t("caronte.last_" + last);
      this.updateQuery();
    },
    setCustomInterval() {
      this.showCustomInterval = true;
      this.customIntervalError = false;

      if (this.customStartDate != null && this.customEndDate != null &&
          moment(this.customStartDate).startOf("day").toDate() > this.customEndDate) {
        // start date is after end date
        this.customIntervalError = true;
      } else if (this.customStartDate != null && this.customEndDate != null) {
        this.filterDate = "custom";
        document.title =
          this.$i18n.t("home.title") +
          ": " +
          moment(this.customStartDate).format("DD MMM YYYY") +
          " - " +
          moment(this.customEndDate).format("DD MMM YYYY");
        this.updateQuery();
      } else {
        // user has to choose start date
        $('#customStartDate').trigger('click')

        setTimeout(function () {
          $('#customStartDate').focus()
        }, 50);
      }
    },
    updateQuery() {
      var query;

      if (this.customStartDate !== null && this.filterDate === 'custom') {
        query = {
          theme: this.lightTheme ? "light" : "dark",
          palette: this.colorPalette,
          last: this.filterDate,
          lang: this.language,
          customStartDate: moment(this.customStartDate).format("YYYY-MM-DD"),
          customEndDate: moment(this.customEndDate).format("YYYY-MM-DD")
        }
      } else {
        query = {
          theme: this.lightTheme ? "light" : "dark",
          palette: this.colorPalette,
          last: this.filterDate,
          lang: this.language
        }
      }
      this.$router.push({ query: query });
    },
    getCurrentPath(route, offset) {
      if (offset) {
        return this.$route.path.split("/")[offset] === route;
      } else {
        return this.$route.path.split("/")[1] === route;
      }
    }
  }
};
</script>

<style>
html {
  height: initial !important;
}
body {
  background: #2d2d2d !important;
  height: initial !important;
}
#app {
  margin: 20px !important;
}

.ui.pointing.menu {
  min-height: 70px !important;
}

.ui.inverted.menu {
  background: #1d1e1e !important;
}
.ui.inverted.segment,
.ui.primary.inverted.segment {
  background: #1d1e1e !important;
}

.datepicker {
  font-size: 1rem;
}

.datepicker input {
  text-align: center;
  border-color: transparent;
  border-radius: 0.28571429rem;
  background-color: #f1f1f1;
  width: 115px;
  padding-left: 10px;
  padding-right: 10px;
  margin-left: 5px;
}
.datepicker.inverted input {
  background-color: #1d1e1e;
  color: white;
  border-color: #5d5e5e;
  border-style: solid;
}

.customIntervalPanel {
  display: flex;
  align-items: center;
}

.mg-left-5 {
  margin-left: 5px !important;
}

.mg-left-10 {
  margin-left: 10px !important;
}

.mg-left-15 {
  margin-left: 15px !important;
}

.datepicker.inverted .vdp-datepicker__calendar {
  background-color: #1d1e1e;
  color: white;
}
.datepicker.inverted .vdp-datepicker__calendar .disabled {
  color: #5d5e5e;
}
.datepicker.inverted .vdp-datepicker__calendar .prev {
  background-color: #5d5e5e;
}
.datepicker.inverted .vdp-datepicker__calendar .next {
  background-color: #5d5e5e;
}
.datepicker.inverted .day__month_btn.up:hover {
  color: #1d1e1e;
}
.datepicker.inverted .month__year_btn.up:hover {
  color: #1d1e1e;
}

.datepicker.customIntervalError input {
  background-color: #ca1010 !important;
  color: white;
}

#customStartDate:focus, #customEndDate:focus, .customDateFocus {
  outline-width: 0;
  border: 1px solid #5d5e5e;
}
</style>
