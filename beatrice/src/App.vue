<template>
  <div id="app">
    <div :class="['ui pointing menu', lightTheme ? '' : 'inverted']">
      <a href="#/" :class="[getCurrentPath('') ? 'active' : '', 'item']">{{$t('home.dashboard')}}</a>
      <a
        href="#/settings"
        :class="[getCurrentPath('settings') ? 'active' : '', 'item']"
      >{{$t('home.settings')}}</a>
      <div class="right menu">
        <div class="item">
          <div class="ui transparent icon input" :class="lightTheme ? '' : 'inverted'">
            <input type="text" placeholder="Search...">
            <i class="search link icon"></i>
          </div>
        </div>
        <div class="item">
          <div
            id="toggleTheme"
            @click="toggleTheme()"
            class="ui button"
            :class="lightTheme ? '' : 'inverted'"
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
      lightTheme: true
    };
  },
  methods: {
    toggleTheme() {
      this.lightTheme = !this.lightTheme;
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
body {
  background: #2d2d2d !important;
}
#app {
  margin: 20px;
}
</style>
