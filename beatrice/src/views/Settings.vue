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
