<template>
  <header class="site-header" :class="{ 'is-scrolled': isScrolled }">
    <div class="header-gold-line"></div>
    <div class="header-container">
      <NuxtLink to="/" class="header-logo">
        <img v-if="siteConfig?.site_logo" :src="siteConfig.site_logo" :alt="siteConfig?.site_name || '北极星移民'" class="logo-img" />
        <template v-else>
          <span class="logo-shield">
            <span class="logo-shield-inner">M</span>
          </span>
          <span class="logo-text">{{ siteConfig?.site_name || '北极星移民' }}</span>
        </template>
      </NuxtLink>

      <nav
        class="header-nav"
        :class="{ 'nav-open': mobileMenuOpen }"
        @touchstart="onNavTouchStart"
        @touchend="onNavTouchEnd"
        @click.self="mobileMenuOpen = false"
      >
        <ul class="nav-list">
          <li
            v-for="(item, index) in navItems"
            :key="item.id"
            class="nav-item"
            :class="{
              'has-children': item.children?.length,
              'is-active': isNavActive(item),
            }"
            @mouseenter="openMega(item.id)"
            @mouseleave="closeMega"
          >
            <span v-if="index > 0" class="nav-separator" aria-hidden="true">·</span>

            <NuxtLink v-if="item.link" :to="item.link" class="nav-link" @click="onNavLinkClick(item, $event)">
              {{ item.label }}
              <span v-if="item.children?.length" class="dropdown-arrow">
                <svg class="chevron-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
              </span>
            </NuxtLink>
            <span v-else class="nav-link" @click="onNavLinkClick(item, $event)">
              {{ item.label }}
              <span v-if="item.children?.length" class="dropdown-arrow">
                <svg class="chevron-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
              </span>
            </span>

            <span v-if="isNavActive(item)" class="nav-active-line"></span>

            <!-- Mobile expand toggle -->
            <button
              v-if="item.children?.length"
              class="mobile-expand-toggle"
              @click.stop="toggleMobileExpand(item.id)"
              :aria-label="mobileExpanded.has(item.id) ? '收起' : '展开'"
            >
              <span :class="{ rotated: mobileExpanded.has(item.id) }">
                <svg class="chevron-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
              </span>
            </button>

            <!-- Mega Panel (desktop) -->
            <div
              v-if="item.children?.length && activeMega === item.id"
              class="mega-panel"
              @mouseenter="openMega(item.id)"
            >
              <div class="mega-arrow"></div>
              <div class="mega-inner">
                <div class="mega-glow-orb"></div>
                <div class="mega-list">
                  <div
                    v-for="child in item.children"
                    :key="child.id"
                    class="mega-card"
                    :class="{ 'has-subs': child.children?.length }"
                  >
                    <NuxtLink v-if="child.link" :to="child.link" class="mega-card-title is-link">
                      <span class="mega-card-title-bar"></span>
                      {{ child.label }}
                    </NuxtLink>
                    <span v-else class="mega-card-title">
                      <span class="mega-card-title-bar"></span>
                      {{ child.label }}
                    </span>

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
                  </div>
                </div>
                <div class="mega-bottom-line"></div>
              </div>
            </div>

            <!-- Mobile: expanded children -->
            <Transition name="mobile-sub">
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
                    <span :class="{ rotated: mobileExpanded.has(child.id) }">
                      <svg class="chevron-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
                    </span>
                  </button>
                </div>

                <Transition name="mobile-sub">
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
                </Transition>
              </li>
            </ul>
            </Transition>
          </li>
        </ul>
      </nav>

      <div class="header-actions">
        <a :href="`tel:${siteConfig?.contact_phone || '400-963-6933'}`" class="header-cta">
          <svg class="header-cta-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
          </svg>
          全国咨询热线：{{ siteConfig?.contact_phone || '400-963-6933' }}
        </a>

        <button
          type="button"
          class="hamburger"
          :class="{ active: mobileMenuOpen }"
          @click.stop="toggleMobileMenu"
          aria-label="Toggle navigation menu"
        >
          <span class="hamburger-bar"></span>
          <span class="hamburger-bar"></span>
          <span class="hamburger-bar"></span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const route = useRoute();
const { siteConfig } = useMygoSiteConfig();
const { navItems, refreshNavigation } = useNavigation();

const mobileMenuOpen = ref(false);
const activeMega = ref<number | null>(null);
const mobileExpanded = ref<Set<number>>(new Set());
const isScrolled = ref(false);

let megaCloseTimer: ReturnType<typeof setTimeout> | null = null;
let navTouchStartX = 0;

onMounted(() => {
  refreshNavigation()
  if (import.meta.client) {
    window.addEventListener('scroll', onScroll, { passive: true });
    isScrolled.value = window.scrollY > 50;
  }
});

onUnmounted(() => {
  if (import.meta.client) {
    window.removeEventListener('scroll', onScroll);
    document.body.style.overflow = '';
  }
});

const onScroll = () => {
  isScrolled.value = window.scrollY > 50;
};

const isNavActive = (item: NavItem): boolean => {
  const path = route.path;
  if (item.link && path === item.link) return true;
  if (item.children?.length) {
    return item.children.some((child: NavItem) => {
      if (child.link && path === child.link) return true;
      return child.children?.some((sub: NavItem) => sub.link && path === sub.link) ?? false;
    });
  }
  return false;
};

interface NavItem {
  id: number;
  label: string;
  link: string | null;
  children: NavItem[];
}

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

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
};

const onNavLinkClick = (item: NavItem, event: MouseEvent) => {
  if (!mobileMenuOpen.value || !item.children?.length) return;
  event.preventDefault();
  toggleMobileExpand(item.id);
};

const onNavTouchStart = (e: TouchEvent) => {
  navTouchStartX = e.touches[0].clientX;
};

const onNavTouchEnd = (e: TouchEvent) => {
  const diff = e.changedTouches[0].clientX - navTouchStartX;
  if (diff > 60) {
    mobileMenuOpen.value = false;
  }
};

watch(mobileMenuOpen, (val) => {
  if (import.meta.client) {
    document.body.style.overflow = val ? 'hidden' : '';
  }
  if (!val) {
    mobileExpanded.value = new Set();
  }
});

watch(() => route.path, () => {
  mobileMenuOpen.value = false;
});
</script>

<style scoped>
/* ==================== Header Base ==================== */

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
  transition: height var(--duration-slow) var(--ease-out),
              background var(--duration-slow) var(--ease-out),
              box-shadow var(--duration-slow) var(--ease-out);
}

.site-header.is-scrolled {
  height: var(--header-scrolled-height);
  background: rgba(10, 22, 40, 0.98);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.35);
  border-bottom-color: rgba(200, 150, 62, 0.18);
}

/* Top gold line */
.header-gold-line {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent 5%, var(--accent) 50%, transparent 95%);
  opacity: 0.7;
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

/* ==================== Logo ==================== */

.header-logo {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-shield {
  width: 34px;
  height: 34px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  border-radius: 6px;
  transform: rotate(45deg);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 0 14px rgba(200, 150, 62, 0.25);
}

.logo-shield-inner {
  transform: rotate(-45deg);
  color: #fff;
  font-weight: 800;
  font-size: 15px;
  line-height: 1;
}

.logo-text {
  font-size: 19px;
  font-weight: 700;
  color: var(--bg-white);
  letter-spacing: 2px;
  line-height: 1.2;
}

.logo-img {
  height: 32px;
  width: auto;
}

/* ==================== Nav ==================== */

.header-nav {
  display: flex;
}

.nav-list {
  display: flex;
  align-items: center;
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
}

.nav-separator {
  color: rgba(200, 150, 62, 0.25);
  font-size: 13px;
  margin: 0 1px;
  user-select: none;
  pointer-events: none;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  color: rgba(255, 255, 255, 0.78);
  font-size: 14px;
  font-weight: 500;
  border-radius: 6px;
  transition: all var(--duration-fast) var(--ease-out);
  white-space: nowrap;
  cursor: default;
  position: relative;
}

a.nav-link {
  cursor: pointer;
}

.nav-link:hover {
  color: #fff;
  background-color: rgba(200, 150, 62, 0.1);
}

.nav-item.is-active .nav-link {
  color: #fff;
}

/* Active indicator underline */
.nav-active-line {
  position: absolute;
  bottom: -1px;
  left: 50%;
  width: calc(100% - 24px);
  height: 2px;
  background: var(--accent);
  border-radius: 1px;
  transform: translateX(-50%) scaleX(0);
  animation: activeLineIn var(--duration-normal) var(--ease-spring) forwards;
}

@keyframes activeLineIn {
  to {
    transform: translateX(-50%) scaleX(1);
  }
}

.dropdown-arrow {
  display: inline-flex;
  align-items: center;
  opacity: 0.5;
  transition: opacity var(--duration-fast) var(--ease-out);
}

.chevron-icon {
  display: block;
}

/* ==================== Mega Panel ==================== */

.mega-panel {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  padding-top: 12px;
  min-width: 520px;
  max-width: 680px;
  animation: megaIn var(--duration-normal) var(--ease-out);
  z-index: 101;
}

@keyframes megaIn {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-6px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

/* Arrow pointer */
.mega-arrow {
  position: absolute;
  top: 4px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 7px solid transparent;
  border-right: 7px solid transparent;
  border-bottom: 7px solid var(--accent);
  filter: drop-shadow(0 -1px 2px rgba(0,0,0,0.06));
}

.mega-inner {
  position: relative;
  background-color: #FAFBFC;
  border-radius: var(--radius-xl);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12), 0 0 0 1px rgba(200, 150, 62, 0.06);
  padding: 16px;
  overflow: hidden;
}

/* Decorative glow orb */
.mega-glow-orb {
  position: absolute;
  top: -40px;
  right: -40px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.06) 0%, transparent 70%);
  pointer-events: none;
}

.mega-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  position: relative;
  z-index: 1;
}

/* Card-style group */
.mega-card {
  background: #fff;
  border: 1px solid rgba(200, 150, 62, 0.08);
  border-radius: 10px;
  padding: 10px 12px 8px;
  transition: border-color var(--duration-fast) var(--ease-out),
              box-shadow var(--duration-fast) var(--ease-out);
}

.mega-card:hover {
  border-color: rgba(200, 150, 62, 0.2);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
}

.mega-card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  padding: 4px 0 6px;
  transition: color var(--duration-fast) var(--ease-out);
}

.mega-card-title-bar {
  width: 3px;
  height: 14px;
  background: rgba(200, 150, 62, 0.3);
  border-radius: 2px;
  flex-shrink: 0;
  transition: background var(--duration-fast) var(--ease-out);
}

.mega-card-title.is-link {
  cursor: pointer;
}

.mega-card-title.is-link:hover {
  color: var(--primary);
}

.mega-card-title.is-link:hover .mega-card-title-bar {
  background: var(--accent);
}

.mega-card.has-subs .mega-card-title {
  border-bottom: 1px solid rgba(200, 150, 62, 0.08);
  padding-bottom: 8px;
  margin-bottom: 4px;
}

.mega-subitems {
  padding-top: 2px;
}

.mega-subitem {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: var(--text-secondary);
  padding: 6px 10px 6px 11px;
  border-radius: 6px;
  transition: all var(--duration-fast) var(--ease-out);
  cursor: default;
}

.mega-subitem.is-link {
  cursor: pointer;
}

.mega-subitem.is-link:hover {
  background-color: rgba(200, 150, 62, 0.06);
  color: var(--primary);
  transform: translateX(4px);
}

.mega-subitem-dot {
  width: 4px;
  height: 4px;
  background: rgba(200, 150, 62, 0.4);
  border-radius: 50%;
  margin-right: 8px;
  flex-shrink: 0;
  transition: all var(--duration-fast) var(--ease-out);
}

.mega-subitem.is-link:hover .mega-subitem-dot {
  background: var(--accent);
  width: 6px;
  height: 6px;
}

/* Bottom decorative line */
.mega-bottom-line {
  height: 1px;
  margin-top: 14px;
  background: linear-gradient(90deg, transparent, rgba(200, 150, 62, 0.15), transparent);
}

/* ==================== Header Actions (CTA + Hamburger) ==================== */

.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

/* ==================== CTA Button ==================== */

.header-cta {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 20px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  border-radius: 6px;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
  white-space: nowrap;
  letter-spacing: 0.3px;
  box-shadow: 0 2px 12px rgba(200, 150, 62, 0.25);
  transition: all var(--duration-fast) var(--ease-out);
  flex-shrink: 0;
}

.header-cta-icon {
  width: 15px;
  height: 15px;
  flex-shrink: 0;
}

.header-cta:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(200, 150, 62, 0.4);
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
  height: 2.5px;
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
  padding: 4px 8px;
  cursor: pointer;
  transition: color var(--duration-fast) var(--ease-out);
}

.mobile-expand-toggle:hover {
  color: #fff;
}

.mobile-expand-toggle span {
  display: inline-flex;
  align-items: center;
  transition: transform var(--duration-normal) var(--ease-out);
}

.mobile-expand-toggle span.rotated {
  transform: rotate(180deg);
}

/* ==================== Mobile ==================== */

@media (max-width: 767px) {
  .logo-text {
    font-size: 16px;
    letter-spacing: 1px;
  }

  .logo-shield {
    width: 28px;
    height: 28px;
  }

  .logo-shield-inner {
    font-size: 13px;
  }

  .nav-separator {
    display: none;
  }

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
    width: 100%;
    height: 100vh;
    height: 100dvh;
    background: rgba(10, 20, 37, 0.92);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    padding-top: 80px;
    transform: translateX(100%);
    transition: transform 0.35s cubic-bezier(0.4, 0, 0.2, 1);
    z-index: 100;
    overflow-y: auto;
    overscroll-behavior: contain;
  }

  .header-nav.nav-open {
    transform: translateX(0);
  }

  /* Mobile gold accent line at top of nav overlay */
  .header-nav.nav-open::before {
    content: '';
    position: absolute;
    top: 0;
    left: 10%;
    right: 10%;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(200, 150, 62, 0.4), transparent);
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
    flex-wrap: wrap;
    justify-content: center;
  }

  .nav-item.has-children {
    display: flex;
    cursor: pointer;
  }

  .nav-link {
    font-size: 17px;
    font-weight: 500;
    padding: 14px 20px 14px 44px;
    justify-content: flex-start;
    width: 100%;
    border-radius: 0;
  }

  .nav-link .dropdown-arrow {
    display: none;
  }

  .nav-item.is-active .nav-link {
    color: var(--accent);
  }

  .nav-active-line {
    display: none;
  }

  .mega-panel {
    display: none;
  }

  .mobile-expand-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    left: 0;
    top: 14px;
    width: 44px;
    height: 44px;
  }

  .mobile-submenu {
    width: 100%;
    padding: 4px 0 8px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 0 0 8px 8px;
    border-left: 1px solid rgba(200, 150, 62, 0.15);
    margin-left: 12px;
    overflow: hidden;
  }

  /* Transition: mobile submenu expand/collapse */
  .mobile-sub-enter-active,
  .mobile-sub-leave-active {
    transition: max-height 0.35s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.25s ease, padding 0.3s ease;
    overflow: hidden;
  }

  .mobile-sub-enter-from,
  .mobile-sub-leave-to {
    max-height: 0 !important;
    opacity: 0;
    padding-top: 0 !important;
    padding-bottom: 0 !important;
  }

  .mobile-sub-enter-to,
  .mobile-sub-leave-from {
    max-height: 800px;
    opacity: 1;
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
    padding: 10px 16px 10px 36px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 15px;
    text-align: left;
    width: 100%;
    transition: color 0.15s ease;
  }

  .mobile-subitem-label.is-link:active {
    color: var(--bg-white);
  }

  .mobile-expand-toggle.sub {
    position: absolute;
    left: 4px;
    top: 6px;
    width: 32px;
    height: 32px;
  }

  .mobile-submenu.l3 {
    padding: 0 0 6px 20px;
    margin-left: 0;
    background: rgba(255, 255, 255, 0.015);
    border-left: 1px solid rgba(200, 150, 62, 0.08);
    width: 100%;
  }

  .mobile-subitem.l3 {
    border-bottom: none;
  }

  .mobile-subitem-label.l3 {
    font-size: 14px;
    padding: 8px 16px 8px 44px;
    color: rgba(255, 255, 255, 0.5);
    text-align: left;
  }

  .mobile-subitem-label.l3.is-link:active {
    color: rgba(255, 255, 255, 0.75);
  }

  .header-cta {
    font-size: 0;
    padding: 6px;
    gap: 0;
    border-radius: 6px;
    width: 34px;
    height: 34px;
    justify-content: center;
    margin-left: 6px;
    background: transparent;
    box-shadow: none;
    color: var(--bg-white);
  }

  .header-cta:hover {
    transform: none;
    box-shadow: none;
    color: var(--accent-light);
  }

  .header-cta-icon {
    width: 20px;
    height: 20px;
  }
}
</style>
