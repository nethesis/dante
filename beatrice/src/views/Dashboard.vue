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
  <div class="ui">
    <h1
      class="ui header adjust-header"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >{{$t('home.dashboard')}}</h1>

    <div
      v-show="view.isLoading"
      class="ui active dimmer"
      :class="$parent.lightTheme ? 'inverted' : ''"
    >
      <div class="ui text loader">Loading</div>
    </div>

    <button
      v-show="!view.isLoading && gridLayout.length > 0"
      @click="toggleMode()"
      class="compact ui button right floated grey"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >{{mode == 'edit' ? $t('dashboard.edit_done') : $t('dashboard.edit_widgets')}}</button>

    <button
      v-show="!view.isLoading && gridLayout.length > 0"
      @click="openAddElement()"
      class="ui compact labeled icon button right floated blue"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >
      <i class="add icon"></i>
      {{$t('dashboard.add_widget')}}
    </button>
    <button
      v-show="!view.isLoading && gridLayout.length > 0 && mode == 'edit'"
      @click="resetLayout()"
      class="ui compact labeled icon button right floated red"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >
      <i class="redo icon"></i>
      {{$t('dashboard.reset_default')}}
    </button>

    <div
      v-show="!view.isLoading && gridLayout.length == 0"
      class="ui placeholder segment"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >
      <div class="ui icon header">
        <i class="cube icon"></i>
        {{$t('dashboard.no_widgets')}}
      </div>
      <button
        @click="openAddElement()"
        class="ui compact labeled icon button right floated blue"
        :disabled="mode == 'edit'"
        :class="$parent.lightTheme ? '' : 'inverted'"
      >
        <i class="add icon"></i>
        {{$t('dashboard.add_widget')}}
      </button>
    </div>

    <grid-layout
      v-if="!view.isLoading"
      :layout.sync="gridLayout"
      :col-num="12"
      :row-height="20"
      :is-draggable="mode == 'edit' ? true : false"
      :is-resizable="false"
      :is-mirrored="false"
      :vertical-compact="true"
      :margin="[5, 5]"
      :use-css-transforms="true"
      @layout-updated="layoutUpdated"
    >
      <grid-item
        v-for="(item,k) in gridLayout"
        :key="k"
        :x="item.x"
        :y="item.y"
        :w="item.w"
        :h="item.h"
        :i="item.i"
        :class="[item.type == 'chart' ? '' : 'empty', item.highlight ? 'highlight' : $parent.searchString.length > 0 ? 'lowlight' : '', mode == 'edit' ? $parent.lightTheme ? 'on-edit' : 'on-edit-dark' : '']"
        :isResizable="mode == 'edit' ? true : false"
        @resized="itemResized"
      >
        <!-- CLOSE BUTTON -->
        <button
          v-if="mode == 'edit'"
          class="ui compact icon button red mini adjust-close-icon"
          :class="$parent.lightTheme ? '' : 'inverted'"
          @click="removeElement(item)"
        >
          <i class="remove icon adjust-remove"></i>
        </button>
        <!-- END CLOSE BUTTON -->

        <!-- CHART -->
        <div
          v-if="item.type == 'chart' && item.data && item.data.series && item.data.series.length == 0"
          class="ui active dimmer"
          :class="$parent.lightTheme ? 'inverted' : ''"
        >
          <div class="ui indeterminate text loader">{{$t('dashboard.retrieving_data')}}</div>
        </div>
        <div
          class="ui statistics"
          v-if="item.type == 'chart'  && item.data && item.data.series && item.data.series.length > 0"
        >
          <div class="statistic">
            <div class="text value">
              <Chart
                :chartId="item.id"
                :type="item.data.type"
                :series="item.data.series"
                :categories="item.data.categories"
                :width="item.width"
                :height="item.height"
                :theme="$parent.lightTheme"
                :palette="$parent.colorPalette"
                :title="item.data.title"
                :unit="item.data.unit"
                :labels="item.data.labels"
                :class="mode == 'edit' ? 'adjust-content' : ''"
              />
            </div>
          </div>
        </div>
        <!-- END CHART -->

        <!-- COUNTER -->
        <div
          v-if="item.type == 'counter' && item.data && (item.data.series.length == 0)"
          class="ui active dimmer"
          :class="$parent.lightTheme ? 'inverted' : ''"
        >
          <div class="ui indeterminate text loader">{{$t('dashboard.retrieving_data')}}</div>
        </div>
        <div
          v-if="item.type == 'counter' && item.data && (item.data.series.length > 0)"
          class="ui three statistics"
          :class="[$parent.lightTheme ? '' : 'inverted', 'mini']"
        >
          <div class="statistic">
            <div class="label adjust-label-counter">{{item.data.title || '-'}}</div>
            <div class="value">{{item.data.value || 0 | formatter(item.data.unit)}}</div>
          </div>
          <div class="statistic">
            <div class="label adjust-label-counter"></div>
            <div class="value">
              <Chart
                :chartId="item.id"
                type="area"
                :series="item.data.series"
                :categories="item.data.categories"
                :width="item.width"
                :height="item.height"
                :theme="$parent.lightTheme"
                :palette="$parent.colorPalette"
                :title="item.data.title"
                :labels="item.data.labels"
                :sparkline="true"
                :unit="item.data.unit"
              />
            </div>
          </div>
          <div class="statistic">
            <div class="label adjust-label-counter">
              {{$t('dashboard.trend')}}
              <br>
              <span class="ui header grey">({{$t('home.last_'+$parent.filterDate)}})</span>
            </div>
            <div
              class="ui value header"
              :class="item.data.trend <= 0 ? 'red' : 'green'"
            >{{item.data.trend > 0 ? '+' : ''}}{{item.data.trend | formatter(item.data.trendType)}}</div>
          </div>
        </div>
        <!-- END COUNTER -->

        <!-- TABLE -->
        <div
          v-if="item.type == 'table' && item.data && item.data.rows && item.data.rows.length == 0"
          class="ui active dimmer"
          :class="$parent.lightTheme ? 'inverted' : ''"
        >
          <div class="ui indeterminate text loader">{{$t('dashboard.retrieving_data')}}</div>
        </div>
        <span
          v-if="item.type == 'table' && item.data && item.data.rows && item.data.rows.length > 0"
          class="ui header"
          :class="$parent.lightTheme ? '' : 'inverted'"
        >
          <h5 class="adjust-title-table">{{item.data.title.toUpperCase()}}</h5>
        </span>
        <table
          v-if="item.type == 'table' && item.data && item.data.rows && item.data.rows.length > 0"
          class="ui striped selectable table"
          :class="[$parent.lightTheme ? '' : 'inverted', item.data.rowHeader ? 'definition' : '']"
        >
          <thead>
            <tr>
              <th
                v-for="(h,hk) in item.data.columnHeader"
                :key="hk"
              >{{h.length > 0 ? $t(item.id+'.'+h) : h}}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(r,rk) in item.data.rows" :key="rk">
              <td v-for="(i,ik) in r" :key="ik">{{i | formatter(item.data.unit)}}</td>
            </tr>
          </tbody>
        </table>
        <!-- END TABLE -->

        <!-- LABEL -->
        <div
          v-if="item.type == 'label' && item.data && !item.data.value"
          class="ui active dimmer"
          :class="$parent.lightTheme ? 'inverted' : ''"
        >
          <div class="ui indeterminate text loader">{{$t('dashboard.retrieving_data')}}</div>
        </div>
        <div
          v-if="item.type == 'label' && item.data && item.data.value"
          class="ui one statistics full-box"
          :class="[$parent.lightTheme ? '' : 'inverted', 'mini']"
        >
          <div class="statistic full-box">
            <div class="label adjust-label-counter">{{item.data.title || '-'}}</div>
            <div class="value">{{item.data.value || 0}}</div>
          </div>
        </div>
        <!-- END LABEL -->

        <!-- LIST -->
        <div
          v-if="item.type == 'list' && item.data && !item.data.list"
          class="ui active dimmer"
          :class="$parent.lightTheme ? 'inverted' : ''"
        >
          <div class="ui indeterminate text loader">{{$t('dashboard.retrieving_data')}}</div>
        </div>
        <span
          v-if="item.type == 'list' && item.data && item.data.list"
          class="ui header"
          :class="$parent.lightTheme ? '' : 'inverted'"
        >
          <h5 class="adjust-title-table">{{item.data.title.toUpperCase()}}</h5>
        </span>
        <div
          v-if="item.type == 'list' && item.data && item.data.list"
          class="ui middle aligned divided list large ordered selection adjust-list"
          :class="$parent.lightTheme ? '' : 'inverted'"
        >
          <div v-for="(l,lk) in item.data.list" :key="lk" class="item">
            <div class="content">
              <div class="ui header" :class="setListTitle(lk)">{{l.name}}</div>
              {{l.count | formatter(item.data.unit)}}
            </div>
          </div>
        </div>
        <!-- END LIST -->

        <!-- TITLE -->
        <span
          v-if="item.type == 'title'"
          class="ui header"
          :class="$parent.lightTheme ? '' : 'inverted'"
          @dblclick="mode == 'edit' ? editTitle(item, true) : undefined"
        >
          <h2 v-show="!item.isEdit" class="title-pad">{{item.title}}</h2>

          <div
            v-show="item.isEdit"
            class="ui transparent action input mini adjust-input-container"
            :class="$parent.lightTheme ? '' : 'inverted'"
          >
            <input
              v-model="item.newTitle"
              autofocus
              type="text"
              class="ui input massive adjust-input"
              :placeholder="$t('dashboard.insert_new_title')"
            >
            <i
              @click="item.newTitle && item.newTitle.length > 0 ? editTitle(item, false) : null"
              class="icon adjust-icon"
              :class="[item.newTitle && item.newTitle.length > 0 ? '' : 'disabled', item.newTitle.length > 0 ? 'green check' : 'grey ban']"
            ></i>
          </div>
        </span>
        <!-- END TITLE -->
      </grid-item>
    </grid-layout>

    <div class="ui small modal">
      <div class="header">{{$t('dashboard.add_widget')}}</div>
      <div class="content">
        <div class="ui five column grid link cards">
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'chart' ? 'add-widget-selected' : '']"
              @click="setNewElement('chart')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="chart area icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.chart')}}</a>
              </div>
            </div>
          </div>
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'counter' ? 'add-widget-selected' : '']"
              @click="setNewElement('counter')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="percent icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.counter')}}</a>
              </div>
            </div>
          </div>
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'table' ? 'add-widget-selected' : '']"
              @click="setNewElement('table')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="table icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.table')}}</a>
              </div>
            </div>
          </div>
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'list' ? 'add-widget-selected' : '']"
              @click="setNewElement('list')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="trophy icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.list')}}</a>
              </div>
            </div>
          </div>
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'label' ? 'add-widget-selected' : '']"
              @click="setNewElement('label')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="font icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.label')}}</a>
              </div>
            </div>
          </div>
          <div
            v-if="newObject.selected == 'chart' || newObject.selected == 'counter' || newObject.selected == 'table' || newObject.selected == 'list' || newObject.selected == 'label'"
            class="ui form big grid row centered vertical segment"
          >
            <div class="inline fields">
              <label>{{$t('dashboard.choose_' + newObject.selected)}}</label>
              <div class="field">
                <select v-model="newObject.widget" class="ui inline dropdown">
                  <option
                    v-for="(c,ck) in freeWidgets[newObject.selected]"
                    :key="ck"
                    :value="c"
                  >{{$t(c.name+'.title')}}</option>
                </select>
              </div>
            </div>
          </div>
        </div>
        <div class="ui divider"></div>
        <div class="ui one column grid link cards adjust-card">
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'title' ? 'add-widget-selected' : '']"
              @click="setNewElement('title')"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="heading icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.title')}}</a>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="actions">
        <button @click="closeModal()" class="ui red cancel button">
          <i class="remove icon"></i>
          {{$t('dashboard.cancel')}}
        </button>
        <button
          :disabled="newObject.selected.length == 0 || !newObject.widget"
          @click="addElement()"
          class="ui green ok button"
        >
          <i class="checkmark icon"></i>
          {{$t('dashboard.add')}}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import VueGridLayout from "vue-grid-layout";
import Chart from "@/components/Chart.vue";

var moment = require("moment");

export default {
  name: "dashboard",
  components: {
    Chart,
    GridLayout: VueGridLayout.GridLayout,
    GridItem: VueGridLayout.GridItem
  },
  watch: {
    "$parent.searchString": function(search) {
      this.gridLayout.map(function(g) {
        if (
          search.length > 0 &&
          (g.type == "chart" ||
            g.type == "counter" ||
            g.type == "label" ||
            g.type == "table") &&
          JSON.stringify(g)
            .toLowerCase()
            .includes(search.toLowerCase())
        ) {
          g.highlight = true;
        } else {
          g.highlight = false;
        }
      });
      this.$forceUpdate();
    },
    "$parent.filterDate": function() {
      this.getLayout();
    }
  },
  mounted() {
    var context = this;
    this.getLayout();
    this.getWidgets();
  },
  data() {
    var offset = 30;
    var widgetDefaults = {
      chart: {
        w: 6,
        h: 14,
        width: window.innerWidth / 2 - (offset + 3.5 * 6),
        height: window.innerHeight / 3
      },
      counter: {
        w: 6,
        h: 4,
        width: window.innerWidth / 6,
        height: window.innerHeight / 12
      },
      table: {
        w: 6,
        h: 6
      },
      list: {
        w: 6,
        h: 6
      },
      label: {
        w: 6,
        h: 4
      },
      title: {
        w: 6,
        h: 2
      }
    };

    return {
      offset: offset,
      mode: "view",
      widgetDefaults: widgetDefaults,
      gridLayout: [],
      freeWidgets: [],
      newObject: this.initNewObject(),
      view: {
        isLoading: true
      },
      apiHost: this.$root.$options.apiHost
    };
  },
  methods: {
    initNewObject() {
      return {
        selected: "",
        widget: null
      };
    },
    toggleMode(mode) {
      this.mode = this.mode == "edit" ? "view" : "edit";

      if (this.mode == "view") {
        this.setLayout(this.gridLayout);
      }
    },
    editTitle(item, edit) {
      if (edit) {
        item.newTitle = item.title ? item.title : "";
        item.isEdit = true;
        this.$forceUpdate();
      } else {
        item.title = item.newTitle;
        item.isEdit = false;
        this.$forceUpdate();
      }
    },

    setListTitle(index) {
      switch (index) {
        case 0:
          return "huge";
        case 1:
          return "large";
        case 2:
          return "medium";
        case 3:
          return "small";
        default:
          return "tiny";
      }
    },
    closeModal() {
      $(".small.modal").modal("hide");
    },

    layoutUpdated: function(newLayout) {
      this.setLayout(newLayout);
    },
    itemResized: function(i, newH, newW, newHPx, newWPx) {
      var defaultW = 0;
      var defaultH = 0;

      switch (this.gridLayout[i].type) {
        case "chart":
          defaultW =
            window.innerWidth / (12 / newW) - (this.offset + 7.5 * newW);
          defaultH = window.innerHeight / 3.5 / (12 / newH);
          break;
        case "counter":
          defaultW =
            window.innerWidth / (12 / newW) / 3 - (newW == 12 ? 10 : 0);
          defaultH = window.innerHeight / 3.5 / (12 / newH);
      }

      this.gridLayout[i].width = defaultW;
      this.gridLayout[i].height = defaultH;
    },

    openAddElement() {
      this.newObject = this.initNewObject();
      $(".small.modal").modal("show");
    },
    setNewElement(element) {
      this.newObject.selected = element;
      this.newObject.widget = element == "title" ? { name: "title" } : null;
      this.$forceUpdate();
    },
    addElement() {
      var type = this.newObject.selected;

      // prepare new element position
      this.gridLayout.map(function(elem) {
        elem.y++;
        elem.i++;

        return elem;
      });

      // prepare new object structure
      var obj = this.widgetDefaults[type];
      obj.x = 0;
      obj.y = 0;
      obj.i = 0;
      obj.type = type;
      obj.id = this.newObject.widget.name;

      obj.data = {
        series: []
      };

      obj.newTitle = "";
      obj.title = this.$i18n.t("dashboard.empty_" + obj.type + "_title");

      // add element to grid
      this.gridLayout.unshift(JSON.parse(JSON.stringify(obj)));

      // get widget info
      if (type != "title") {
        this.getWidgetData(obj.id, obj.i);
      }

      // save layout
      this.setLayout(this.gridLayout);

      // close modal
      $(".small.modal").modal("hide");
    },
    removeElement(item) {
      this.gridLayout.splice(this.gridLayout.indexOf(item), 1);
      this.gridLayout.map(function(item, index) {
        item.i = index;
      });

      this.mode = this.gridLayout.length == 0 ? "view" : "edit";
      this.setLayout(this.gridLayout);
    },
    getWidgets() {
      this.$http.get(this.apiHost + "/miners").then(
        success => {
          for (var w in success.body.miners) {
            var widget = success.body.miners[w];

            if (this.freeWidgets[widget.type]) {
              this.freeWidgets[widget.type].push(widget);
            } else {
              this.freeWidgets[widget.type] = [widget];
            }
          }
        },
        error => {
          console.error(error);
        }
      );
    },
    getLayout() {
      this.view.isLoading = true;
      this.$http.get(this.apiHost + "/layout").then(
        success => {
          var layouts = success.body.layout;

          for (var l in layouts) {
            var layout = layouts[l];

            layout.w = success.body.default
              ? this.widgetDefaults[layout.type].w
              : layout.w;
            layout.h = success.body.default
              ? this.widgetDefaults[layout.type].h
              : layout.h;
            layout.width = success.body.default
              ? this.widgetDefaults[layout.type].width
              : (window.innerWidth * layout.width) / 100;
            layout.height = success.body.default
              ? this.widgetDefaults[layout.type].height
              : (window.innerHeight * layout.height) / 100;
            layout.data = {
              series: []
            };
            layout.newTitle = "";
            layout.title = layout.text;

            if (layout.type != "title") {
              this.getWidgetData(layout.id, l);
            }
          }

          this.gridLayout = layouts;
          this.view.isLoading = false;
        },
        error => {
          this.gridLayout = [];
          this.view.isLoading = false;
          console.error(error);
        }
      );
    },
    getWidgetData(widget, index) {
      var startDate = "";
      var endDate = moment();

      switch (this.$parent.filterDate) {
        case "week":
          startDate = moment()
            .subtract(7, "days")
            .startOf("day");
          break;

        case "month":
          startDate = moment()
            .subtract(1, "months")
            .startOf("day");
          break;

        case "halfyear":
          startDate = moment()
            .subtract(6, "months")
            .startOf("day");
          break;
      }

      this.$http
        .get(
          this.apiHost +
            "/widget/" +
            widget +
            "?startDate=" +
            startDate.format("YYYY-MM-DD") +
            "&endDate=" +
            endDate.format("YYYY-MM-DD")
        )
        .then(
          success => {
            // get body data
            var widget = success.body.widget;

            if (widget) {
              this.gridLayout[index].data.title =
                this.$i18n.t(widget.minerId + ".title") || "-";
              this.gridLayout[index].data.type = widget.chartType || "line";
              this.gridLayout[index].data.labels = widget.labels || [];
              this.gridLayout[index].data.value = widget.value || 0;
              this.gridLayout[index].data.series =
                widget.series || widget.trendSeries || [];
              this.gridLayout[index].data.categories =
                widget.categories || widget.trendCategories || [];
              this.gridLayout[index].data.list = widget.data || [];
              this.gridLayout[index].data.rows = widget.rows || [];
              this.gridLayout[index].data.rowHeader = widget.rowHeader || false;
              this.gridLayout[index].data.columnHeader =
                widget.columnHeader || [];
              this.gridLayout[index].data.trend = widget.trend;
              this.gridLayout[index].data.trendType = widget.trendType;
              this.gridLayout[index].data.unit = widget.unit;
              this.gridLayout[index].data.aggregationType =
                widget.aggregationType;

              // calculate correct sizes
              if (widget.type == "list") {
                this.gridLayout[index].h =
                  (widget.data.length < 10 ? 7 : 10) + 1.5 * widget.data.length;
              }
              if (widget.type == "table") {
                this.gridLayout[index].h = 6 + 1.5 * widget.rows.length;
              }
            }

            this.$forceUpdate();
          },
          error => {
            console.error(error);
          }
        );
    },
    setLayout(newLayout) {
      newLayout = newLayout.map(function(w) {
        return {
          id: w.id,
          x: w.x,
          y: w.y,
          i: w.i,
          type: w.type,
          w: w.w,
          h: w.h,
          width: (w.width * 100) / window.innerWidth,
          height: (w.height * 100) / window.innerHeight,
          text: w.title || ""
        };
      });
      this.$http.post(this.apiHost + "/layout", { layout: newLayout }).then(
        success => {
          console.info("saved");
        },
        error => {
          console.error(error);
        }
      );
    },
    resetLayout() {
      this.$http.delete(this.apiHost + "/layout").then(
        success => {
          console.info("reset");
          this.getLayout();
        },
        error => {
          console.error(error);
        }
      );
    }
  }
};
</script>

<style>
.on-edit {
  border-radius: 0.28571429rem !important;
  outline: 2px dashed #e0e1e2 !important;
}
.on-edit-dark {
  border-radius: 0.28571429rem !important;
  outline: 2px dashed #828282 !important;
}

.vue-grid-item.vue-grid-placeholder {
  background: #e0e1e2 !important;
  border: 2px solid black !important;
  border-radius: 0.28571429rem;
}

.vue-grid-item.vue-draggable-dragging {
  border: 2px dashed #e0e1e2 !important;
  border-radius: 0.28571429rem !important;
}

.empty {
  border-radius: 0.28571429rem !important;
}

.title-pad {
  padding: 3px !important;
}

.adjust-input-container {
  padding: 3px !important;
}
.adjust-input:focus {
  border-right: none !important;
}
.adjust-icon {
  margin-top: 5px !important;
}
.adjust-icon:hover {
  cursor: pointer !important;
}
.adjust-header {
  display: inline-block !important;
}
.adjust-remove {
  cursor: pointer !important;
}
.adjust-content {
  display: inline-block !important;
}
.adjust-label-counter {
  margin-bottom: 15px !important;
}
.adjust-image-icon {
  padding: 10px !important;
  cursor: pointer !important;
}
.adjust-close-icon {
  position: absolute !important;
  right: -2px !important;
  top: 1px !important;
  z-index: 99999 !important;
}
.adjust-title-table {
  padding-left: 10px !important;
  padding-top: 2px !important;
}
.adjust-list {
  margin-top: 5px !important;
}
.adjust-card {
  margin-top: -12px !important;
}

.highlight {
  border-radius: 0.28571429rem !important;
  outline: 2px solid #54c8ff !important;
}
.lowlight {
  opacity: 0.25 !important;
}

.ui.segment.inverted
  > .ui
  > .vue-grid-layout
  > .vue-grid-item
  > .vue-resizable-handle {
  filter: invert(100%) !important;
}
.ui.segment > .ui > .vue-grid-layout > .vue-grid-item > .vue-resizable-handle {
  filter: invert(0%) !important;
}

.add-widget-selected {
  border: 2px solid #54c8ff !important;
  margin-top: -2px !important;
}

.ui.dimmer {
  background-color: rgb(29, 30, 30, 0.75) !important;
}
.ui.inverted.dimmer {
  background-color: rgba(255, 255, 255, 0.75) !important;
}

.ui.statistics {
  align-items: center !important;
  margin-top: 0px !important;
  margin-left: 0px !important;
  margin-right: 0px !important;
  margin-bottom: 0px !important;
}
.ui.ui.statistics > .statistic {
  margin-top: 0px !important;
  margin-left: 0px !important;
  margin-right: 0px !important;
  margin-bottom: 0px !important;
  max-width: 33.333333% !important;
}
.vue-resizable-handle {
  z-index: 2 !important;
}

.ui.inverted.definition.table tfoot:not(.full-width) th:first-child,
.ui.inverted.definition.table thead:not(.full-width) th:first-child {
  background: #1d1e1e !important;
  box-shadow: -1px -1px 0 1px #1d1e1e !important;
}

.full-box {
  width: 100% !important;
  height: 100% !important;
}

.ui.list .list > .item > .content,
.ui.list > .item > .content {
  margin-left: 10px !important;
}
</style>