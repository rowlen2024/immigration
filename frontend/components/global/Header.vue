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
            :key="item.label"
            class="nav-item"
            @mouseenter="openDropdown(item.label)"
            @mouseleave="closeDropdown"
          >
            <NuxtLink :to="item.link" class="nav-link">
              {{ item.label }}
              <span v-if="item.children && item.children.length" class="dropdown-arrow">&#9662;</span>
            </NuxtLink>
            <ul
              v-if="item.children.length && activeDropdown === item.label"
              class="dropdown-menu"
            >
              <li v-for="child in item.children" :key="child.label">
                <NuxtLink :to="child.link" class="dropdown-item">
                  {{ child.label }}
                </NuxtLink>
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
const activeDropdown = ref<string | null>(null);

onMounted(() => {
  fetchSiteConfig();
});

const openDropdown = (label: string) => {
  activeDropdown.value = label;
};

const closeDropdown = () => {
  activeDropdown.value = null;
};

watch(mobileMenuOpen, (val) => {
  if (import.meta.client) {
    document.body.style.overflow = val ? 'hidden' : '';
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

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  padding-top: 6px;
  min-width: 180px;
  background-color: var(--bg-white);
  border-radius: 10px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  padding-bottom: 6px;
  animation: dropdownIn 0.2s ease-out;
}

@keyframes dropdownIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-item {
  display: block;
  padding: 9px 18px;
  color: var(--text-secondary);
  font-size: 13px;
  transition: all 0.15s ease;
}

.dropdown-item:hover {
  background-color: var(--bg-light);
  color: var(--primary);
}

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
    justify-content: center;
    transform: translateX(100%);
    transition: transform 0.3s ease;
    z-index: 100;
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

  .nav-link {
    font-size: 17px;
    font-weight: 500;
    padding: 14px 0;
    justify-content: center;
    width: 100%;
    border-radius: 0;
  }

  .dropdown-menu {
    position: static;
    background-color: transparent;
    box-shadow: none;
    text-align: center;
    padding: 0 0 8px;
  }

  .dropdown-item {
    color: rgba(255, 255, 255, 0.65);
    padding: 7px 16px;
    font-size: 14px;
  }

  .dropdown-item:hover {
    background-color: transparent;
    color: var(--bg-white);
  }

  .header-cta {
    display: none;
  }
}
</style>
