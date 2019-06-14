<template>
  <div class="ui">
    <h1 class="ui header" :class="$parent.lightTheme ? '' : 'inverted'">{{$t('home.dashboard')}}</h1>
    <div class="ui two column grid">
      <Chart
        id="chart-one"
        type="line"
        :series="series1"
        :width="width/2-offset"
        :height="height"
        :theme="$parent.lightTheme"
      />
      <Chart
        id="chart-two"
        type="area"
        :series="series2"
        :width="width/2-offset"
        :height="height"
        :theme="$parent.lightTheme"
      />
    </div>
  </div>
</template>

<script>
import Chart from "@/components/Chart.vue";

export default {
  name: "dashboard",
  components: {
    Chart
  },
  data() {
    return {
      series1: [
        {
          name: "High - 2013",
          data: [28, 29, 33, 36, 32, 32, 33]
        },
        {
          name: "Low - 2013",
          data: [12, 11, 14, 18, 17, 13, 13]
        }
      ],
      series2: [
        {
          name: "High - 2013",
          data: [28, 29, 33, 36, 32, 32, 33]
        },
        {
          name: "Low - 2013",
          data: [12, 11, 14, 18, 17, 13, 13]
        }
      ],
      width: window.innerWidth,
      height: window.innerHeight / 3,
      offset: 60
    };
  },
  methods: {}
};
</script>

<style>
</style>