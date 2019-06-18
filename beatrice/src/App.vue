<template>
  <div id="app">
    <div :class="['ui pointing menu', lightTheme ? '' : 'inverted']">
      <a href="#/" :class="[getCurrentPath('') ? 'active' : '', 'item']">{{$t('home.dashboard')}}</a>
      <div class="item">
        <div class="ui buttons" :class="lightTheme ? '' : 'inverted'">
          <div class="ui button" :class="lightTheme ? '' : 'inverted'">{{$t('home.last_day')}}</div>
          <div class="ui button" :class="lightTheme ? '' : 'inverted'">{{$t('home.last_week')}}</div>
          <div
            class="ui button active"
            :class="lightTheme ? '' : 'inverted'"
          >{{$t('home.last_month')}}</div>
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
              >Palette 1</div>
              <div
                @click="setPalette(2)"
                class="item"
                :class="colorPalette == 'palette2' ? 'selected' : ''"
              >Palette 2</div>
              <div
                @click="setPalette(3)"
                class="item"
                :class="colorPalette == 'palette3' ? 'selected' : ''"
              >Palette 3</div>
              <div
                @click="setPalette(4)"
                class="item"
                :class="colorPalette == 'palette4' ? 'selected' : ''"
              >Palette 4</div>
              <div
                @click="setPalette(5)"
                class="item"
                :class="colorPalette == 'palette5' ? 'selected' : ''"
              >Palette 5</div>
              <div
                @click="setPalette(6)"
                class="item"
                :class="colorPalette == 'palette6' ? 'selected' : ''"
              >Palette 6</div>
              <div
                @click="setPalette(7)"
                class="item"
                :class="colorPalette == 'palette7' ? 'selected' : ''"
              >Palette 7</div>
              <div
                @click="setPalette(8)"
                class="item"
                :class="colorPalette == 'palette8' ? 'selected' : ''"
              >Palette 8</div>
              <div
                @click="setPalette(9)"
                class="item"
                :class="colorPalette == 'palette9' ? 'selected' : ''"
              >Palette 9</div>
            </div>
          </div>
        </div>
      </div>
      <div class="right menu">
        <div class="item">
          <div class="ui transparent icon input" :class="lightTheme ? '' : 'inverted'">
            <input v-model="searchString" type="text" placeholder="Search...">
            <i class="search link icon"></i>
          </div>
        </div>
        <div class="item">
          <div
            id="toggleTheme"
            @click="toggleTheme()"
            class="ui button"
            :class="lightTheme ? 'black' : 'inverted'"
          >{{lightTheme ? $t('settings.dark_theme') : $t('settings.light_theme')}}</div>
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
      lightTheme: this.$route.query.theme == "light",
      colorPalette: this.$route.query.palette || "palette1",
      searchString: ""
    };
  },
  methods: {
    toggleTheme() {
      this.lightTheme = !this.lightTheme;
      this.updateQuery();
    },
    setPalette(paletteNumber) {
      this.colorPalette = "palette" + paletteNumber;
      this.updateQuery();
    },
    updateQuery() {
      this.$router.push({
        query: {
          theme: this.lightTheme ? "light" : "dark",
          palette: this.colorPalette
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

.ui.inverted.menu {
  background: #1d1e1e !important;
}
.ui.inverted.segment,
.ui.primary.inverted.segment {
  background: #1d1e1e !important;
}
</style>
