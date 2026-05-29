<template>
  <section class="consult-cta">
    <h3>{{ title }}</h3>
    <p>{{ description }}</p>
    <NuxtLink v-if="buttonLink" :to="buttonLink" class="btn-primary">{{ buttonText }}</NuxtLink>
  </section>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  description: string
  buttonText?: string
  buttonLink?: string
}>(), {
  buttonText: '免费咨询',
  buttonLink: '/contact',
})
</script>

<style scoped>
.consult-cta {
  text-align: center;
  padding: 48px 0;
  background-color: var(--bg-light);
  border-radius: var(--radius-lg);
  margin: 40px 0;
}

.consult-cta h3 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.consult-cta p {
  font-size: 15px;
  color: var(--text-secondary);
  margin-bottom: 24px;
}

@media (max-width: 767px) {
  .consult-cta {
    padding: 36px 20px;
    margin: 32px 0;
    border-radius: var(--radius-md);
  }

  .consult-cta h3 {
    font-size: 20px;
  }

  .consult-cta .btn-primary {
    width: 100%;
    padding: 14px 0;
    font-size: 16px;
  }
}
</style>
