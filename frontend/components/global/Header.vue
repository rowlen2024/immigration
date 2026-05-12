<template>
  <header class="site-header">
    <div class="header-container">
      <NuxtLink to="/" class="header-logo">
        <span v-if="!siteConfig?.site_logo" class="logo-mark">M</span>
        <img v-if="siteConfig?.site_logo" :src="siteConfig.site_logo" :alt="siteConfig?.site_name || '北极星移民'" class="logo-img" />
        <span v-else class="logo-text">{{ siteConfig?.site_name || '北极星移民' }}</span>
      </NuxtLink>

      <nav class="header-nav" :class="{ 'nav-open': mobileMenuOpen }">
        <ul class="nav-list">
          <li
            v-for="item in navItems"
            :key="item.id"
            class="nav-item"
            :class="{ 'has-children': item.children?.length }"
            @mouseenter="openMega(item.id)"
            @mouseleave="closeMega"
          >
            <!-- Level 1 link -->
            <NuxtLink v-if="item.link" :to="item.link" class="nav-link">
              {{ item.label }}
              <span v-if="item.children?.length" class="dropdown-arrow">&#9662;</span>
            </NuxtLink>
            <span v-else class="nav-link">
              {{ item.label }}
              <span v-if="item.children?.length" class="dropdown-arrow">&#9662;</span>
            </span>

            <!-- Mobile: expand toggle -->
            <button
              v-if="item.children?.length"
              class="mobile-expand-toggle"
              @click.stop="toggleMobileExpand(item.id)"
              :aria-label="mobileExpanded.has(item.id) ? '收起' : '展开'"
            >
              <span :class="{ rotated: mobileExpanded.has(item.id) }">&#9662;</span>
            </button>

            <!-- Mega Panel (desktop) -->
            <div
              v-if="item.children?.length && activeMega === item.id"
              class="mega-panel"
              @mouseenter="openMega(item.id)"
            >
              <ul class="mega-list">
                <li
                  v-for="child in item.children"
                  :key="child.id"
                  class="mega-group"
                  :class="{ 'has-subs': child.children?.length }"
                >
                  <!-- Level 2 title -->
                  <NuxtLink v-if="child.link" :to="child.link" class="mega-group-title is-link">
                    {{ child.label }}
                  </NuxtLink>
                  <span v-else class="mega-group-title">{{ child.label }}</span>

                  <!-- Level 3 items -->
                  <ul v-if="child.children?.length" class="mega-subitems">
                    <li v-for="sub in child.children" :key="sub.id">
                      <NuxtLink v-if="sub.link" :to="sub.link" class="mega-subitem is-link">
                        <span class="mega-subitem-dot"></span>{{ sub.label }}
                      </NuxtLink>
                      <span v-else class="mega-subitem">
                        <span class="mega-subitem-dot"></span>{{ sub.label }}
                      </span>
                    </li>
                  </ul>
                </li>
              </ul>
            </div>

            <!-- Mobile: expanded children (accordion) -->
            <ul
              v-if="item.children?.length && mobileExpanded.has(item.id)"
              class="mobile-submenu"
            >
              <li
                v-for="child in item.children"
                :key="child.id"
                class="mobile-subitem"
                :class="{ 'has-subs': child.children?.length }"
              >
                <div class="mobile-subitem-row">
                  <NuxtLink
                    v-if="child.link"
                    :to="child.link"
                    class="mobile-subitem-label is-link"
                    @click="mobileMenuOpen = false"
                  >
                    {{ child.label }}
                  </NuxtLink>
                  <span v-else class="mobile-subitem-label">{{ child.label }}</span>

                  <button
                    v-if="child.children?.length"
                    class="mobile-expand-toggle sub"
                    @click.stop="toggleMobileExpand(child.id)"
                    :aria-label="mobileExpanded.has(child.id) ? '收起' : '展开'"
                  >
                    <span :class="{ rotated: mobileExpanded.has(child.id) }">&#9662;</span>
                  </button>
                </div>

                <!-- Level 3 mobile submenu -->
                <ul
                  v-if="child.children?.length && mobileExpanded.has(child.id)"
                  class="mobile-submenu l3"
                >
                  <li v-for="sub in child.children" :key="sub.id" class="mobile-subitem l3">
                    <NuxtLink
                      v-if="sub.link"
                      :to="sub.link"
                      class="mobile-subitem-label l3 is-link"
                      @click="mobileMenuOpen = false"
                    >
                      {{ sub.label }}
                    </NuxtLink>
                    <span v-else class="mobile-subitem-label l3">{{ sub.label }}</span>
                  </li>
                </ul>
              </li>
            </ul>
          </li>
        </ul>
      </nav>

      <NuxtLink to="/contact" class="header-cta">免费咨询</NuxtLink>

      <button
        class="hamburger"
        :class="{ active: mobileMenuOpen }"
        @click="mobileMenuOpen = !mobileMenuOpen"
        aria-label="Toggle navigation menu"
      >
        <span class="hamburger-bar"></span>
        <span class="hamburger-bar"></span>
        <span class="hamburger-bar"></span>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();
const { navItems, fetchNav } = useNavigation();

const mobileMenuOpen = ref(false);
const activeMega = ref<number | null>(null);
const mobileExpanded = ref<Set<number>>(new Set());

let megaCloseTimer: ReturnType<typeof setTimeout> | null = null;

onMounted(() => {
  fetchSiteConfig();
});

const openMega = (id: number) => {
  if (megaCloseTimer) {
    clearTimeout(megaCloseTimer);
    megaCloseTimer = null;
  }
  activeMega.value = id;
};

const closeMega = () => {
  megaCloseTimer = setTimeout(() => {
    activeMega.value = null;
  }, 150);
};

const toggleMobileExpand = (id: number) => {
  const next = new Set(mobileExpanded.value);
  if (next.has(id)) {
    next.delete(id);
  } else {
    next.add(id);
  }
  mobileExpanded.value = next;
};

watch(mobileMenuOpen, (val) => {
  if (import.meta.client) {
    document.body.style.overflow = val ? 'hidden' : '';
  }
  if (!val) {
    mobileExpanded.value = new Set();
  }
});

onUnmounted(() => {
  if (import.meta.client) {
    document.body.style.overflow = '';
  }
});
</script>

<style scoped>
.site-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: var(--header-height);
  background: linear-gradient(180deg, rgba(15, 30, 61, 0.97) 0%, rgba(21, 41, 77, 0.95) 100%);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(200, 150, 62, 0.12);
}

.header-container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 24px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-logo {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--bg-white);
  letter-spacing: -0.5px;
}

.logo-img {
  height: 32px;
  width: auto;
  filter: brightness(10) saturate(0);
}

.logo-mark {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--bg-white);
  font-weight: 800;
  font-size: 14px;
  flex-shrink: 0;
}

.header-nav {
  display: flex;
}

.nav-list {
  display: flex;
  gap: 2px;
}

.nav-item {
  position: relative;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  color: rgba(255, 255, 255, 0.82);
  font-size: 14px;
  font-weight: 500;
  border-radius: 6px;
  transition: all 0.2s ease;
  white-space: nowrap;
  cursor: default;
}

a.nav-link {
  cursor: pointer;
}

.nav-link:hover {
  color: var(--bg-white);
  background-color: rgba(255, 255, 255, 0.08);
}

.dropdown-arrow {
  font-size: 9px;
  opacity: 0.5;
  transition: transform 0.2s ease;
}

/* ==================== Mega Panel (Desktop) ==================== */

.mega-panel {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  padding-top: 6px;
  min-width: 480px;
  max-width: 640px;
  animation: megaIn 0.2s ease-out;
  z-index: 101;
}

@keyframes megaIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.mega-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 24px;
  background-color: var(--bg-white);
  border-radius: 12px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  padding: 20px 24px;
}

.mega-group {
  min-width: 0;
}

.mega-group-title {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  padding: 6px 10px;
  border-radius: 6px;
  transition: all 0.15s ease;
  margin-bottom: 2px;
}

.mega-group-title.is-link {
  cursor: pointer;
}

.mega-group-title.is-link:hover {
  background-color: var(--bg-light);
  color: var(--primary);
}

.mega-group.has-subs .mega-group-title {
  border-bottom: 1px solid var(--bg-light);
  border-radius: 6px 6px 0 0;
  margin-bottom: 0;
  padding-bottom: 8px;
}

.mega-subitems {
  padding: 0 0 4px 10px;
}

.mega-subitem {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: var(--text-secondary);
  padding: 5px 10px;
  border-radius: 4px;
  transition: all 0.15s ease;
  cursor: default;
}

.mega-subitem.is-link {
  cursor: pointer;
}

.mega-subitem.is-link:hover {
  background-color: var(--bg-light);
  color: var(--primary);
}

.mega-subitem-dot {
  margin-right: 6px;
  opacity: 0.35;
  font-size: 10px;
}

/* ==================== CTA Button ==================== */

.header-cta {
  padding: 9px 22px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  color: var(--bg-white);
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(200, 150, 62, 0.2);
  transition: all 0.2s ease;
  white-space: nowrap;
}

.header-cta:hover {
  box-shadow: 0 4px 16px rgba(200, 150, 62, 0.35);
  transform: translateY(-1px);
}

/* ==================== Hamburger ==================== */

.hamburger {
  display: none;
  flex-direction: column;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  z-index: 101;
}

.hamburger-bar {
  width: 24px;
  height: 2px;
  background-color: var(--bg-white);
  border-radius: 1px;
  transition: all 0.3s ease;
}

/* ==================== Mobile expand toggle ==================== */

.mobile-expand-toggle {
  display: none;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.65);
  font-size: 12px;
  padding: 4px 8px;
  cursor: pointer;
  transition: color 0.2s ease;
}

.mobile-expand-toggle:hover {
  color: var(--bg-white);
}

.mobile-expand-toggle span {
  display: inline-block;
  transition: transform 0.25s ease;
}

.mobile-expand-toggle span.rotated {
  transform: rotate(180deg);
}

/* ==================== Mobile: full-screen nav overlay ==================== */

@media (max-width: 767px) {
  .hamburger {
    display: flex;
  }

  .hamburger.active .hamburger-bar:nth-child(1) {
    transform: rotate(45deg) translate(5px, 5px);
  }

  .hamburger.active .hamburger-bar:nth-child(2) {
    opacity: 0;
  }

  .hamburger.active .hamburger-bar:nth-child(3) {
    transform: rotate(-45deg) translate(5px, -5px);
  }

  .header-nav {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(180deg, #0A1425, #0F1E3D, #15294D);
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    padding-top: 80px;
    transform: translateX(100%);
    transition: transform 0.3s ease;
    z-index: 100;
    overflow-y: auto;
  }

  .header-nav.nav-open {
    transform: translateX(0);
  }

  .nav-list {
    flex-direction: column;
    align-items: center;
    gap: 0;
    width: 100%;
    max-width: 320px;
  }

  .nav-item {
    width: 100%;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  }

  .nav-item.has-children {
    display: flex;
    flex-wrap: wrap;
  }

  .nav-link {
    font-size: 17px;
    font-weight: 500;
    padding: 14px 0;
    justify-content: center;
    width: 100%;
    border-radius: 0;
  }

  .mega-panel {
    display: none;
  }

  .mobile-expand-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    right: 0;
    top: 14px;
    width: 44px;
    height: 44px;
  }

  .mobile-submenu {
    width: 100%;
    padding: 0 0 8px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 0 0 8px 8px;
  }

  .mobile-subitem {
    width: 100%;
  }

  .mobile-subitem-row {
    display: flex;
    align-items: center;
    width: 100%;
    position: relative;
  }

  .mobile-subitem-label {
    display: block;
    padding: 10px 32px;
    color: rgba(255, 255, 255, 0.65);
    font-size: 15px;
    text-align: center;
    width: 100%;
    transition: color 0.15s ease;
  }

  .mobile-subitem-label.is-link:active {
    color: var(--bg-white);
  }

  .mobile-expand-toggle.sub {
    position: absolute;
    right: 8px;
    top: 6px;
    width: 32px;
    height: 32px;
  }

  .mobile-submenu.l3 {
    padding: 0 0 6px 16px;
    background: rgba(255, 255, 255, 0.02);
  }

  .mobile-subitem.l3 {
    border-bottom: none;
  }

  .mobile-subitem-label.l3 {
    font-size: 14px;
    padding: 8px 32px;
    color: rgba(255, 255, 255, 0.5);
  }

  .mobile-subitem-label.l3.is-link:active {
    color: rgba(255, 255, 255, 0.8);
  }

  .header-cta {
    display: none;
  }
}
</style>
