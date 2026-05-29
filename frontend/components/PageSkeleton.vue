<template>
  <div class="page-skeleton" :class="`page-skeleton--${variant}`">
    <!-- 首页轮播图 -->
    <template v-if="variant === 'hero'">
      <div class="sk-hero">
        <div class="sk-block sk-hero__badge"></div>
        <div class="sk-block sk-hero__title"></div>
        <div class="sk-block sk-hero__subtitle"></div>
        <div class="sk-hero__actions">
          <div class="sk-block sk-hero__btn"></div>
          <div class="sk-block sk-hero__btn"></div>
        </div>
      </div>
      <div class="sk-trust">
        <div v-for="i in 4" :key="i" class="sk-trust__item">
          <div class="sk-block sk-trust__number"></div>
          <div class="sk-block sk-trust__label"></div>
        </div>
      </div>
      <div class="sk-section">
        <div class="sk-block sk-section__title"></div>
        <div class="sk-block sk-section__subtitle"></div>
        <div class="sk-cards">
          <div v-for="i in 3" :key="i" class="sk-card">
            <div class="sk-block sk-card__image"></div>
            <div class="sk-block sk-card__title"></div>
            <div class="sk-block sk-card__text"></div>
            <div class="sk-block sk-card__text sk-card__text--short"></div>
          </div>
        </div>
      </div>
    </template>

    <!-- 卡片网格（案例/项目列表） -->
    <template v-else-if="variant === 'cards'">
      <div class="sk-page-header">
        <div class="sk-block sk-page-header__title"></div>
        <div class="sk-block sk-page-header__subtitle"></div>
      </div>
      <div class="sk-cards">
        <div v-for="i in count" :key="i" class="sk-card">
          <div class="sk-block sk-card__image"></div>
          <div class="sk-block sk-card__title"></div>
          <div class="sk-block sk-card__text"></div>
          <div class="sk-block sk-card__text sk-card__text--short"></div>
        </div>
      </div>
    </template>

    <!-- 列表（FAQ/咨询） -->
    <template v-else-if="variant === 'list'">
      <div class="sk-page-header">
        <div class="sk-block sk-page-header__title"></div>
        <div class="sk-block sk-page-header__subtitle"></div>
      </div>
      <div class="sk-filters">
        <div v-for="i in 4" :key="i" class="sk-block sk-filter__chip"></div>
      </div>
      <div class="sk-list">
        <div v-for="i in count" :key="i" class="sk-list__item">
          <div class="sk-block sk-list__title"></div>
          <div class="sk-block sk-list__arrow"></div>
        </div>
      </div>
    </template>

    <!-- 详情页（项目/案例详情） -->
    <template v-else-if="variant === 'detail'">
      <div class="sk-detail-hero">
        <div class="sk-block sk-detail-hero__title"></div>
        <div class="sk-block sk-detail-hero__subtitle"></div>
      </div>
      <div class="sk-detail-tabs">
        <div v-for="i in 5" :key="i" class="sk-block sk-detail-tabs__item"></div>
      </div>
      <div class="sk-detail-section" v-for="i in 3" :key="i">
        <div class="sk-block sk-detail-section__title"></div>
        <div class="sk-block sk-detail-section__text"></div>
        <div class="sk-block sk-detail-section__text"></div>
        <div class="sk-block sk-detail-section__text sk-detail-section__text--short"></div>
      </div>
    </template>

    <!-- 文章内容（CMS页面） -->
    <template v-else-if="variant === 'content'">
      <div class="sk-content">
        <div class="sk-block sk-content__title"></div>
        <div class="sk-block sk-content__para"></div>
        <div class="sk-block sk-content__para"></div>
        <div class="sk-block sk-content__para sk-content__para--short"></div>
        <div class="sk-block sk-content__para"></div>
        <div class="sk-block sk-content__para sk-content__para--short"></div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
defineProps({
  variant: {
    type: String as () => 'hero' | 'cards' | 'list' | 'detail' | 'content',
    default: 'cards',
  },
  count: {
    type: Number,
    default: 3,
  },
})
</script>

<style scoped>
.page-skeleton {
  padding: 0;
}

/* ── Shimmer Block ── */
.sk-block {
  background: linear-gradient(
    90deg,
    var(--color-bg-gray, #F0F1F3) 25%,
    #e8e9eb 37%,
    var(--color-bg-gray, #F0F1F3) 63%
  );
  background-size: 400% 100%;
  animation: sk-shimmer 1.6s ease-in-out infinite;
  border-radius: var(--radius-sm, 4px);
}

@keyframes sk-shimmer {
  0% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* ── Page Header (shared) ── */
.sk-page-header {
  text-align: center;
  margin-bottom: 40px;
}

.sk-page-header__title {
  width: 200px;
  height: 32px;
  margin: 0 auto 12px;
}

.sk-page-header__subtitle {
  width: 320px;
  height: 18px;
  margin: 0 auto;
}

/* ── Hero ── */
.sk-hero {
  padding: 80px 20px;
  text-align: center;
  background: var(--color-bg-light, #F8F9FB);
}

.sk-hero__badge {
  width: 140px;
  height: 28px;
  margin: 0 auto 20px;
  border-radius: 20px;
}

.sk-hero__title {
  width: 360px;
  height: 44px;
  margin: 0 auto 16px;
}

.sk-hero__subtitle {
  width: 280px;
  height: 20px;
  margin: 0 auto 32px;
}

.sk-hero__actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.sk-hero__btn {
  width: 140px;
  height: 44px;
  border-radius: var(--radius-md, 8px);
}

/* ── Trust Bar ── */
.sk-trust {
  display: flex;
  justify-content: center;
  gap: 64px;
  padding: 36px 40px;
  background: rgba(15, 27, 45, 0.04);
}

.sk-trust__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.sk-trust__number {
  width: 80px;
  height: 36px;
}

.sk-trust__label {
  width: 60px;
  height: 14px;
}

/* ── Section ── */
.sk-section {
  padding: 60px 20px;
}

.sk-section__title {
  width: 200px;
  height: 32px;
  margin: 0 auto 10px;
}

.sk-section__subtitle {
  width: 280px;
  height: 18px;
  margin: 0 auto 40px;
}

/* ── Cards ── */
.sk-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  max-width: var(--max-width, 1200px);
  margin: 0 auto;
  padding: 0 20px;
}

.sk-card {
  background: var(--color-bg-white, #fff);
  border-radius: var(--radius-lg, 12px);
  overflow: hidden;
  border: 1px solid var(--color-border, #E2E8F0);
}

.sk-card__image {
  height: 200px;
  border-radius: 0;
}

.sk-card__title {
  width: 70%;
  height: 20px;
  margin: 20px 24px 10px;
}

.sk-card__text {
  width: 90%;
  height: 14px;
  margin: 0 24px 8px;
}

.sk-card__text--short {
  width: 55%;
  margin-bottom: 20px;
}

/* ── Filters ── */
.sk-filters {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 32px;
  padding: 0 20px;
}

.sk-filter__chip {
  width: 72px;
  height: 34px;
  border-radius: var(--radius-full, 20px);
}

/* ── List ── */
.sk-list {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 20px;
}

.sk-list__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 20px;
  border-bottom: 1px solid var(--color-border, #E2E8F0);
}

.sk-list__title {
  width: 60%;
  height: 18px;
}

.sk-list__arrow {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

/* ── Detail ── */
.sk-detail-hero {
  padding: 80px 20px;
  background: linear-gradient(135deg, #0F1E3D, #15294D);
  text-align: center;
}

.sk-detail-hero__title {
  width: 300px;
  height: 42px;
  margin: 0 auto 12px;
  background: linear-gradient(90deg,
    rgba(255,255,255,0.08) 25%,
    rgba(255,255,255,0.14) 37%,
    rgba(255,255,255,0.08) 63%
  );
  background-size: 400% 100%;
}

.sk-detail-hero__subtitle {
  width: 400px;
  height: 20px;
  margin: 0 auto;
  background: linear-gradient(90deg,
    rgba(255,255,255,0.06) 25%,
    rgba(255,255,255,0.1) 37%,
    rgba(255,255,255,0.06) 63%
  );
  background-size: 400% 100%;
}

.sk-detail-tabs {
  display: flex;
  gap: 0;
  justify-content: center;
  padding: 14px 20px;
  border-bottom: 1px solid var(--color-border, #E2E8F0);
  background: var(--color-bg-white, #fff);
  position: sticky;
  top: 0;
}

.sk-detail-tabs__item {
  width: 80px;
  height: 18px;
  margin: 0 12px;
}

.sk-detail-section {
  max-width: var(--max-width, 1200px);
  margin: 0 auto;
  padding: 48px 20px;
  border-bottom: 1px solid var(--color-border, #E2E8F0);
}

.sk-detail-section__title {
  width: 180px;
  height: 28px;
  margin-bottom: 24px;
}

.sk-detail-section__text {
  width: 100%;
  height: 16px;
  margin-bottom: 12px;
}

.sk-detail-section__text--short {
  width: 60%;
}

/* ── Content ── */
.sk-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 60px 20px;
}

.sk-content__title {
  width: 260px;
  height: 36px;
  margin-bottom: 32px;
}

.sk-content__para {
  width: 100%;
  height: 16px;
  margin-bottom: 14px;
}

.sk-content__para--short {
  width: 70%;
}

@media (prefers-reduced-motion: reduce) {
  .sk-block {
    animation: none;
    background: var(--color-bg-gray, #F0F1F3);
  }
}

/* ── Responsive ── */
@media (max-width: 1023px) {
  .sk-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .sk-trust {
    gap: 32px;
  }
}

@media (max-width: 767px) {
  .sk-cards {
    grid-template-columns: 1fr;
  }

  .sk-hero__title {
    width: 260px;
    height: 36px;
  }

  .sk-hero__subtitle {
    width: 220px;
  }

  .sk-trust {
    flex-direction: column;
    gap: 16px;
    padding: 24px 20px;
  }

  .sk-detail-hero__title {
    width: 220px;
    height: 32px;
  }

  .sk-detail-hero__subtitle {
    width: 280px;
  }

  .sk-detail-tabs {
    gap: 8px;
    overflow-x: auto;
  }

  .sk-page-header__subtitle {
    width: 240px;
  }
}
</style>
