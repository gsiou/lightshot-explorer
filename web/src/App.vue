<template>
  <div id="app">
    <h1>LightShotExplorer</h1>
    <h3>The weird side of the internet</h3>
    <div class="pure-g">
      <div class="pure-u-1 pure-u-md-1-5"></div>
      <div class="pure-u-1 pure-u-md-3-5">
        <Main :img="$route.params.id"/>
      </div>
      <div class="pure-u-1 pure-u-md-1-5"></div>
    </div>
  </div>
</template>

<script>
import Main from './components/Main.vue'

export default {
  name: 'app',
  components: {
    Main
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
