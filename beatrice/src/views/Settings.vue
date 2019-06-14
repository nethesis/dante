<template>
  <div class="ui">
    <h1 class="ui header" :class="$parent.lightTheme ? '' : 'inverted'">{{$t('home.settings')}}</h1>
  </div>
</template>

<script>
export default {
  name: "settings",
  data() {
    return {};
  }
};
</script>
