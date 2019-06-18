<template>
  <div class="ui">
    <h1
      class="ui header adjust-header"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >{{$t('home.dashboard')}}</h1>

    <button
      v-show="gridLayout.length > 0"
      @click="toggleMode()"
      class="compact ui button right floated grey"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >{{mode == 'edit' ? $t('dashboard.edit_done') : $t('dashboard.edit_widgets')}}</button>

    <button
      v-show="gridLayout.length > 0"
      @click="openAddElement()"
      class="ui compact labeled icon button right floated blue"
      :disabled="mode == 'edit'"
      :class="$parent.lightTheme ? '' : 'inverted'"
    >
      <i class="add icon"></i>
      {{$t('dashboard.add_widget')}}
    </button>

    <div
      v-show="gridLayout.length == 0"
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
        :class="[item.type == 'chart' ? '' : 'empty', item.highlight ? 'highlight' : '', mode == 'edit' ? $parent.lightTheme ? 'on-edit' : 'on-edit-dark' : '']"
        :isResizable="mode == 'edit' ? true : false"
        @resized="itemResized"
      >
        <button
          v-if="mode == 'edit'"
          class="ui compact icon button red mini adjust-close-icon"
          :class="$parent.lightTheme ? '' : 'inverted'"
          @click="removeElement(item)"
        >
          <i class="remove icon adjust-remove"></i>
        </button>
        <div class="ui statistics" v-if="item.type == 'chart'">
          <div class="statistic">
            <div class="text value">
              <Chart
                ref="charts"
                :chartId="item.id"
                type="area"
                :series="series1"
                :width="item.width"
                :height="item.height"
                :theme="$parent.lightTheme"
                :palette="$parent.colorPalette"
                :title="item.title"
                :class="mode == 'edit' ? 'adjust-content' : ''"
              />
            </div>
          </div>
        </div>

        <div
          v-if="item.type == 'counter'"
          class="ui statistics"
          :class="[$parent.lightTheme ? '' : 'inverted', mapTitleSize(item.width)]"
        >
          <div class="statistic">
            <div class="label">{{item.title}}</div>
            <div class="value">{{item.value || 0}}</div>
          </div>
          <div class="statistic">
            <div class="text value">
              <Chart
                ref="charts"
                :chartId="item.id"
                type="line"
                :series="series1"
                :width="item.width"
                :height="item.height"
                :theme="$parent.lightTheme"
                :palette="$parent.colorPalette"
                :title="item.title"
                :sparkline="true"
              />
            </div>
          </div>
        </div>

        <span
          v-if="item.type == 'title'"
          class="ui header"
          :class="$parent.lightTheme ? '' : 'inverted'"
          @dblclick="mode == 'edit' ?  editTitle(item, true) : undefined"
        >
          <h3 v-show="!item.isEdit" class="title-pad">{{item.title}}</h3>

          <div
            v-show="item.isEdit"
            class="ui transparent action input mini"
            :class="$parent.lightTheme ? '' : 'inverted'"
          >
            <input
              v-model="item.newTitle"
              autofocus
              type="text"
              class="adjust-input"
              :placeholder="$t('dashboard.insert_new_title')"
            >
            <i
              @click="item.newTitle && item.newTitle.length > 0 ?  editTitle(item, false) : null"
              class="check icon green adjust-icon"
              :class="[item.newTitle && item.newTitle.length > 0 ? '' : 'disabled']"
            ></i>
          </div>
        </span>
      </grid-item>
    </grid-layout>

    <div class="ui tiny modal">
      <div class="header">{{$t('dashboard.add_widget')}}</div>
      <div class="content">
        <div class="ui three column grid link cards">
          <div class="column">
            <div
              class="ui fluid card"
              :class="[newObject.selected == 'chart' ? 'add-widget-selected' : '']"
              @click="newObject.selected = 'chart'"
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
              @click="newObject.selected = 'counter'"
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
              :class="[newObject.selected == 'title' ? 'add-widget-selected' : '']"
              @click="newObject.selected = 'title'"
            >
              <div class="center aligned image adjust-image-icon">
                <i class="heading icon huge"></i>
              </div>
              <div class="center aligned content">
                <a class="header">{{$t('dashboard.title')}}</a>
              </div>
            </div>
          </div>
          <div
            v-if="newObject.selected == 'chart' || newObject.selected == 'counter'"
            class="ui form big grid row centered vertical segment"
          >
            <div class="inline fields">
              <label>{{newObject.selected == 'chart' ? $t('dashboard.choose_chart') : $t('dashboard.choose_counter')}}</label>
              <div class="field">
                <select class="ui inline dropdown">
                  <option value>Gender</option>
                  <option value="1">Male</option>
                  <option value="0">Female</option>
                </select>
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
          :disabled="newObject.selected.length == 0"
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
          g.type == "chart" &&
          g.tags.toString().includes(search)
        ) {
          g.highlight = true;
        } else {
          g.highlight = false;
        }
      });
      this.$forceUpdate();
    }
  },
  data() {
    var offset = 30;
    var widgetDefaults = {
      chart: {
        w: 6,
        h: 14,
        width: window.innerWidth / 2 - (offset + 7.5 * 6),
        height: window.innerHeight / 3
      },
      counter: {
        w: 6,
        h: 4,
        width: window.innerWidth / 2 - ((window.innerWidth / 2) * 45) / 100,
        height: window.innerHeight / 12
      },
      title: {
        w: 6,
        h: 2
      }
    };
    return {
      series1: [

      ],
      series2: [

      ],
      offset: offset,
      mode: "view",
      widgetDefaults: widgetDefaults,
      gridLayout: [

      ],
      newObject: this.initNewObject()
    };
  },
  methods: {
    initNewObject() {
      return {
        selected: ""
      };
    },
    toggleMode(mode) {
      this.mode = this.mode == "edit" ? "view" : "edit";
    },
    layoutUpdated: function(newLayout) {
      console.log("Updated layout: ", newLayout);
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
            window.innerWidth / (12 / newW) -
            ((window.innerWidth / (12 / newW)) * 45) / 100;
          defaultH = window.innerHeight / 3.5 / (12 / newH);
      }

      this.gridLayout[i].width = defaultW;
      this.gridLayout[i].height = defaultH;
    },
    openAddElement() {
      this.newObject = this.initNewObject();
      $(".tiny.modal").modal("show");
    },
    addElement() {
      var lastY = 0;
      var lastI = 0;
      var type = this.newObject.selected;

      // prepare new element position
      this.gridLayout.map(function(elem) {
        if (elem.y > lastY) lastY = elem.y;
        if (elem.i > lastI) lastI = elem.i;
      });

      var obj = this.widgetDefaults[type];
      obj.x = 0;
      obj.y = lastY + 1;
      obj.i = lastI + 1;
      obj.id = type + (lastI + 1);
      obj.type = type;

      // select widget type
      switch (type) {
        case "chart":
          obj.title = this.$i18n.t("dashboard.empty_chart_title");
          break;
        case "counter":
          obj.title = this.$i18n.t("dashboard.empty_counter_title");
          break;
        case "title":
          obj.title = this.$i18n.t("dashboard.empty_title");
          break;
      }
      this.gridLayout.push(JSON.parse(JSON.stringify(obj)));
      $(".tiny.modal").modal("hide");
    },
    removeElement(item) {
      this.gridLayout.splice(this.gridLayout.indexOf(item), 1);
      this.gridLayout.map(function(item, index) {
        item.i = index;
      });

      this.mode = this.gridLayout.length == 0 ? "view" : "edit";
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
    mapTitleSize(width) {
      switch (true) {
        case width < 40:
          return "mini";
        case width >= 40 && width < 110:
          return "mini";
        case width >= 110 && width < 180:
          return "tiny";
        case width >= 180 && width < 250:
          return "small";
        case width >= 250:
          return "small";
      }
    },
    closeModal() {
      $(".tiny.modal").modal("hide");
    }
  }
};
</script>

<style>
.on-edit {
  border-radius: 0.28571429rem;
  border: 2px dashed #e0e1e2;
  margin-top: -2px;
  margin-left: -2px;
}
.on-edit-dark {
  border-radius: 0.28571429rem;
  border: 2px dashed #828282;
  margin-top: -2px;
  margin-left: -2px;
}
.vue-grid-item.vue-grid-placeholder {
  background: #e0e1e2;
  border: 2px solid black;
  border-radius: 0.28571429rem;
}

.vue-grid-item.vue-draggable-dragging {
  border: 2px dashed #e0e1e2;
  border-radius: 0.28571429rem;
}

.empty {
  border-radius: 0.28571429rem;
}

.title-pad {
  padding: 3px;
}

.adjust-input:focus {
  border-right: none !important;
}
.adjust-icon:hover {
  cursor: pointer;
}
.adjust-header {
  display: inline-block;
}
.adjust-remove {
  cursor: pointer;
}
.adjust-content {
  display: inline-block;
}
.adjust-label-counter {
  margin-bottom: 15px;
}
.adjust-image-icon {
  padding: 10px !important;
  cursor: pointer;
}
.adjust-close-icon {
  position: absolute;
  right: -2px;
  top: 1px;
  z-index: 99;
}

.highlight {
  border-radius: 0.28571429rem;
  border: 2px solid #54c8ff;
  margin-top: -2px;
  margin-left: -2px;
}

.ui.segment.inverted
  > .ui
  > .vue-grid-layout
  > .vue-grid-item
  > .vue-resizable-handle {
  filter: invert(100%);
}
.ui.segment > .ui > .vue-grid-layout > .vue-grid-item > .vue-resizable-handle {
  filter: invert(0%);
}

.add-widget-selected {
  border: 2px solid #54c8ff !important;
  margin-top: -2px !important;
}

.ui.statistics {
  display: table !important;
  table-layout: fixed;
  margin-top: 0px !important;
  margin-left: 0px !important;
  margin-right: 0px !important;
  margin-bottom: 0px !important;
  width: 100%;
}
.ui.statistics > .statistic {
  display: table-cell !important;
  vertical-align: middle;
}
.vue-resizable-handle {
  z-index: 2;
}
</style>