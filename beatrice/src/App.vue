<template>
  <div id="app">
    <div :class="['ui pointing menu', lightTheme ? '' : 'inverted']">
      <a href="#/" :class="[getCurrentPath('') ? 'active' : '', 'item']">{{$t('home.title')}}</a>
      <div class="item">
        <div class="ui buttons" :class="lightTheme ? '' : 'inverted'">
          <button
            @click="setFilterDate('day')"
            class="ui button"
            :class="[lightTheme ? '' : 'inverted', filterDate == 'day' ? 'active' : '']"
          >{{$t('home.last_day')}}</button>
          <button
            @click="setFilterDate('week')"
            class="ui button"
            :class="[lightTheme ? '' : 'inverted', filterDate == 'week' ? 'active' : '']"
          >{{$t('home.last_week')}}</button>
          <button
            @click="setFilterDate('month')"
            class="ui button"
            :class="[lightTheme ? '' : 'inverted', filterDate == 'month' ? 'active' : '']"
          >{{$t('home.last_month')}}</button>
        </div>
      </div>
      <div class="item">
        <div class="ui compact menu" :class="lightTheme ? '' : 'inverted'">
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
      <div class="right menu">
        <div class="item">
          <div class="ui transparent icon input" :class="lightTheme ? '' : 'inverted'">
            <input v-model="searchString" type="text" placeholder="Search...">
            <i
              @click="searchString.length > 0 ? searchString = '' : undefined"
              :class="[searchString.length > 0 ? 'remove link' : 'search', 'icon']"
            ></i>
          </div>
        </div>
        <div class="item">
          <div
            id="toggleTheme"
            @click="setTheme()"
            class="ui button"
            :class="lightTheme ? 'black' : 'inverted'"
          >{{lightTheme ? $t('home.dark_theme') : $t('home.light_theme')}}</div>
        </div>
      </div>
    </div>
    <div :class="['ui segment', lightTheme ? '' : 'inverted']">
      <router-view/>
    </div>
  </div>
</template>

<script>
export default {
  name: "home",
  data() {
    return {
      lightTheme: this.$route.query.theme
        ? this.$route.query.theme == "light"
        : "light",
      colorPalette: this.$route.query.palette || "palette1",
      filterDate: this.$route.query.last || "day",
      language: this.$route.query.lang || "en",
      searchString: ""
    };
  },
  methods: {
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
      this.updateQuery();
    },
    updateQuery() {
      this.$router.push({
        query: {
          theme: this.lightTheme ? "light" : "dark",
          palette: this.colorPalette,
          last: this.filterDate,
          lang: this.language
        }
      });
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
  margin: 20px;
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
</style>
