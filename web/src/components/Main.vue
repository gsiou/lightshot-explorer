<template>
  <div>
    <div v-if="img">
      <h1 v-if="loading">Loading...</h1>
      <img v-if="!loading" :src="image"/>
      <br/>
      <button class="pure-button" :disabled="this.loading" v-on:click="getPrev">Prev</button>
      <button class="pure-button pure-button-primary" :disabled="this.loading" v-on:click="getNext">Next</button>
    </div>
    <div v-else>
      <p>Load an image to get started</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Main',
  props: {
    img: String
  },
  data: function () {
    return {
      image: '',
      next: '',
      prev: '',
      loading: false
    }
  },
  mounted () {
    if (this.img) {
      this.getImage(this.img)
    }
  },
  methods: {
    getImage(image) {
      this.$router.push({ path: '/' + image })
      this.loading = true
      window.fetch('http://localhost:12345/image/' + image)
      .then(res => res.json())
      .then(res => {
        console.log(res)
        this.image = res.image
        this.next = res.next
        this.prev = res.prev
        this.loading = false
      })  
    },
    getNext() {
      this.getImage(this.next)
    },
    getPrev() {
      this.getImage(this.prev)
    }
  }
}
</script>

<style scoped>
img {
  max-height: 80vh;
  width: 100%;
}
</style>
