<template>
  <form class="simc-input" @submit.prevent="submitSimc" aria-labelledby="simc-input-label">
    <label for="simc-input-textarea" id="simc-input-label">SimC Input</label>
    <textarea
      id="simc-input-textarea"
      v-model="simcText"
      placeholder="Paste your SimC input here..."
      rows="10"
      cols="50"
      required
      aria-required="true"
    ></textarea>
    <button type="submit">Submit</button>
  </form>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const simcText = ref('')
const submitSimc = () => {
  // TODO
  console.log('SimC Input:', simcText.value)
}
</script>
<style scoped>
form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 600px;
  margin: auto;
}
</style>
