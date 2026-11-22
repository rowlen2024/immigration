<template>
  <div v-if="hasContent" class="contact-sidebar">
    <!-- 联系电话 -->
    <a v-if="siteConfig?.contact_phone_2" class="cs-item cs-item-link" :href="`tel:${siteConfig.contact_phone_2}`" @click="showMobilePhone(siteConfig.contact_phone_2!, '联系电话', '工作日 9:00-18:00', $event)">
      <span class="cs-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
      </span>
      <span class="cs-label">联系电话</span>
      <div class="cs-tooltip">
        <span class="tt-phone">{{ siteConfig.contact_phone_2 }}</span>
        <div class="tt-desc">工作日 9:00-18:00</div>
        <span class="tt-copy" :class="{ copied: copied === siteConfig.contact_phone_2 }" @click.stop.prevent="copy(siteConfig.contact_phone_2!)">{{ copied === siteConfig.contact_phone_2 ? '已复制' : '复制号码' }}</span>
      </div>
    </a>

    <!-- 客服电话 -->
    <a v-if="siteConfig?.contact_phone" class="cs-item cs-item-link" :href="`tel:${siteConfig.contact_phone}`" @click="showMobilePhone(siteConfig.contact_phone!, '客服电话', '24小时咨询热线', $event)">
      <span class="cs-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 18v-6a9 9 0 0 1 18 0v6"/><path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z"/></svg>
      </span>
      <span class="cs-label">客服电话</span>
      <div class="cs-tooltip">
        <span class="tt-phone">{{ siteConfig.contact_phone }}</span>
        <div class="tt-desc">24小时咨询热线</div>
        <span class="tt-copy" :class="{ copied: copied === siteConfig.contact_phone }" @click.stop.prevent="copy(siteConfig.contact_phone!)">{{ copied === siteConfig.contact_phone ? '已复制' : '复制号码' }}</span>
      </div>
    </a>

    <!-- 微信 -->
    <button v-if="siteConfig?.contact_wechat" type="button" class="cs-item" aria-haspopup="dialog" @click="showMobileQR('wechat', $event)">
      <span class="cs-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>
      </span>
      <span class="cs-label">微信</span>
      <span class="cs-tooltip">
        <span class="tt-name">微信扫码咨询</span>
        <img :src="siteConfig.contact_wechat" alt="微信二维码" class="tt-qr-img" loading="lazy" />
      </span>
    </button>

    <!-- 公众号 -->
    <button v-if="siteConfig?.contact_wechat_mp" type="button" class="cs-item" aria-haspopup="dialog" @click="showMobileQR('mp', $event)">
      <span class="cs-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="2" width="16" height="20" rx="2"/><line x1="8" y1="6" x2="16" y2="6"/><line x1="8" y1="10" x2="16" y2="10"/><line x1="8" y1="14" x2="12" y2="14"/></svg>
      </span>
      <span class="cs-label">公众号</span>
      <span class="cs-tooltip">
        <span class="tt-name">关注微信公众号</span>
        <span class="tt-desc">获取最新移民资讯</span>
        <img :src="siteConfig.contact_wechat_mp" alt="微信公众号二维码" class="tt-qr-img" loading="lazy" />
      </span>
    </button>

    <!-- 视频号 -->
    <button v-if="siteConfig?.contact_wechat_video" type="button" class="cs-item" aria-haspopup="dialog" @click="showMobileQR('video', $event)">
      <span class="cs-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2"/></svg>
      </span>
      <span class="cs-label">视频号</span>
      <span class="cs-tooltip">
        <span class="tt-name">关注企业视频号</span>
        <span class="tt-desc">了解更多移民动态</span>
        <img :src="siteConfig.contact_wechat_video" alt="企业视频号二维码" class="tt-qr-img" loading="lazy" />
      </span>
    </button>

    <PublicDialog :id="contactDialogId" v-model:open="contactDialogOpen" :title="contactPopup?.title ?? ''">
      <div v-if="contactPopup" class="contact-popup-content">
        <img v-if="contactPopup.type === 'qr'" :src="contactPopup.image" :alt="contactPopup.title" class="mobile-qr-img" loading="lazy" />
        <div v-else class="mobile-phone-number">{{ contactPopup.number }}</div>
        <div v-if="contactPopup.description" class="mobile-qr-desc">{{ contactPopup.description }}</div>
        <div v-if="contactPopup.type === 'phone'" class="mobile-phone-actions">
          <a :href="`tel:${contactPopup.number}`" class="public-button public-button--gold">立即拨打</a>
          <button type="button" class="public-button" :class="{ copied: copied === contactPopup.number }" @click="contactPopup?.type === 'phone' && copy(contactPopup.number)">{{ copied === contactPopup.number ? '已复制' : '复制号码' }}</button>
        </div>
        <p v-if="contactPopup.type === 'phone' && copyFailed" class="mobile-qr-desc" role="status">复制失败，请手动复制号码。</p>
      </div>
    </PublicDialog>

    <!-- 回到顶部 -->
    <button type="button" class="cs-scroll-top" :class="{ visible: showScrollTop }" aria-label="回到顶部" @click="scrollToTop">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="18 15 12 9 6 15"/></svg>
    </button>
  </div>
</template>

<script setup lang="ts">
const { siteConfig } = useMygoSiteConfig();

const showScrollTop = ref(false);
const copied = ref<string | null>(null);
const copyFailed = ref(false);
let copyTimer: ReturnType<typeof setTimeout> | null = null;

type ContactPopup = {
  title: string;
  description?: string;
} & ({ type: 'qr'; image: string } | { type: 'phone'; number: string });

const contactDialogId = useId();
const contactDialogOpen = ref(false);
const contactPopup = ref<ContactPopup | null>(null);

function openContactDialog(popup: ContactPopup, event: MouseEvent) {
  (event.currentTarget as HTMLElement).focus({ preventScroll: true });
  contactPopup.value = popup;
  copyFailed.value = false;
  contactDialogOpen.value = true;
}

const showMobileQR = (type: 'wechat' | 'mp' | 'video', event: MouseEvent) => {
  const c = siteConfig.value;
  if (!c) return;
  const channels = {
    wechat: { image: c.contact_wechat, title: '微信扫码咨询' },
    mp: { image: c.contact_wechat_mp, title: '关注微信公众号', description: '获取最新移民资讯' },
    video: { image: c.contact_wechat_video, title: '关注企业视频号', description: '了解更多移民动态' },
  };
  const { image, ...details } = channels[type];
  if (image) openContactDialog({ type: 'qr', image, ...details }, event);
};

const showMobilePhone = (number: string, label: string, desc: string, event: MouseEvent) => {
  if (window.innerWidth > 890) return;
  event.preventDefault();
  openContactDialog({ type: 'phone', number, title: label, description: desc }, event);
};

const hasContent = computed(() => {
  const c = siteConfig.value;
  return !!(c?.contact_phone || c?.contact_phone_2 || c?.contact_wechat || c?.contact_wechat_mp || c?.contact_wechat_video);
});

async function copy(text: string) {
  copyFailed.value = false;
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    const focused = document.activeElement;
    const ta = document.createElement('textarea');
    ta.value = text;
    ta.readOnly = true;
    ta.style.position = 'fixed';
    ta.style.left = '-9999px';
    ta.style.top = '-9999px';
    const host = focused?.closest('dialog[open]') ?? document.body;
    host.appendChild(ta);
    let success = false;
    try {
      ta.focus({ preventScroll: true });
      ta.select();
      success = document.activeElement === ta && document.execCommand('copy');
    } catch {
      success = false;
    } finally {
      ta.remove();
      if (focused instanceof HTMLElement) focused.focus({ preventScroll: true });
    }
    if (!success) {
      copied.value = null;
      copyFailed.value = true;
      return;
    }
  }
  copied.value = text;
  if (copyTimer) clearTimeout(copyTimer);
  copyTimer = setTimeout(() => { copied.value = null; }, 1500);
}

function scrollToTop() {
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  window.scrollTo({ top: 0, behavior: reduceMotion ? 'instant' : 'smooth' });
}

function onScroll() {
  showScrollTop.value = window.scrollY > 600;
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll);
  if (copyTimer) clearTimeout(copyTimer);
});
</script>

<style scoped>
.contact-sidebar {
  position: fixed;
  right: 24px;
  bottom: 94px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.cs-item {
  position: relative;
  width: 66px;
  height: 66px;
  margin-top: 8px;
  background: #fff;
  border: 1px solid var(--accent);
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-family: inherit;
  transition: background-color var(--public-motion-fast) var(--ease-out);
}

.cs-item:hover {
  background: var(--accent);
}

.cs-item:hover .cs-icon,
.cs-item:hover .cs-label {
  color: #fff;
}

.cs-icon {
  line-height: 1;
  color: var(--accent);
  transition: color var(--public-motion-fast) var(--ease-out);
  display: flex;
  align-items: center;
  justify-content: center;
}

.cs-label {
  font-size: 11px;
  color: var(--accent);
  margin-top: 3px;
  transition: color var(--public-motion-fast) var(--ease-out);
  font-weight: 500;
}

.cs-item::before {
  content: '';
  position: absolute;
  left: -14px;
  top: -4px;
  bottom: -4px;
  width: 14px;
}

/* Tooltip */
.cs-tooltip {
  visibility: hidden;
  opacity: 0;
  position: absolute;
  right: calc(100% + 8px);
  top: 50%;
  transform: translateY(-50%);
  width: 240px;
  background: #fff;
  border-radius: 10px;
  padding: 14px 10px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, .12);
  text-align: center;
  z-index: 1001;
  transition: opacity var(--public-motion-fast) var(--ease-in),
              visibility 0s var(--public-motion-fast);
}

.cs-tooltip::after {
  content: '';
  position: absolute;
  right: -10px;
  top: 50%;
  transform: translateY(-50%);
  border: 5px solid transparent;
  border-left: 6px solid #fff;
}

.cs-item:hover .cs-tooltip,
.cs-item:focus-visible .cs-tooltip {
  visibility: visible;
  opacity: 1;
  transition: opacity var(--public-motion-fast) var(--ease-out),
              visibility 0s;
}

.tt-phone {
  font-size: 18px;
  font-weight: 700;
  color: #15294D;
  margin-bottom: 4px;
  display: block;
  text-decoration: none;
}

.tt-desc {
  display: block;
  font-size: 12px;
  color: #8e8ea0;
  margin-bottom: 8px;
}

.tt-copy {
  display: inline-block;
  padding: 4px 14px;
  font-size: 12px;
  color: var(--accent);
  border: 1px solid var(--accent);
  border-radius: 4px;
  cursor: pointer;
  transition: color var(--public-motion-fast) var(--ease-out),
              background-color var(--public-motion-fast) var(--ease-out);
}

.tt-copy:hover {
  background: var(--accent);
  color: #fff;
}

.tt-copy.copied {
  background: var(--accent);
  color: #fff;
}

.tt-name {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #15294D;
  margin-bottom: 4px;
}

.tt-qr-img {
  display: block;
  width: 130px;
  aspect-ratio: 1;
  border-radius: 6px;
  margin: 8px auto 0;
  border: 1px solid #e8e8ef;
  object-fit: contain;
  background: #f8f9fb;
}

/* 回到顶部 */
.cs-scroll-top {
  display: none;
  width: 66px;
  height: 66px;
  margin-top: 8px;
  background: #fff;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  color: #8e8ea0;
  transition: color var(--public-motion-fast) var(--ease-out),
              background-color var(--public-motion-fast) var(--ease-out),
              border-color var(--public-motion-fast) var(--ease-out);
}

.cs-scroll-top.visible {
  display: flex;
}

.cs-scroll-top:hover {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

/* 移动端 */
@media (max-width: 890px) {
  .contact-sidebar {
    flex-direction: row;
    right: 0;
    left: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-top: 1px solid var(--border-color);
    z-index: 1001;
    padding-bottom: env(safe-area-inset-bottom, 0px);
  }

  .cs-item {
    flex: 1;
    margin-top: 0;
    border: none !important;
    border-radius: 0;
    width: auto;
    height: auto;
    min-height: 48px;
    padding: 10px 4px;
    background: transparent;
    flex-direction: row;
    gap: 6px;
  }

  .cs-item-link {
    text-decoration: none;
  }

  .cs-item:hover {
    background: transparent !important;
  }

  .cs-icon {
    color: #15294D !important;
  }

  .cs-icon svg {
    width: 18px;
    height: 18px;
  }

  .cs-label {
    font-size: 12px;
    color: #15294D !important;
    font-weight: 500;
    margin-top: 0;
    white-space: nowrap;
  }

  .cs-tooltip {
    display: none !important;
  }

  .cs-scroll-top {
    display: flex !important;
    position: fixed;
    right: 12px;
    bottom: calc(68px + env(safe-area-inset-bottom, 0px));
    width: 40px;
    height: 40px;
    background: rgba(255, 255, 255, 0.92);
    border: 1px solid var(--border-color);
    border-radius: 50%;
    box-shadow: var(--shadow-md);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    color: var(--color-text-secondary);
    z-index: 1002;
    margin-top: 0;
  }

  .cs-scroll-top.visible {
    display: flex !important;
  }
}

.contact-popup-content {
  text-align: center;
}

.mobile-qr-desc {
  font-size: 13px;
  color: #8e8ea0;
  margin-top: 8px;
}

.mobile-qr-img {
  width: 200px;
  aspect-ratio: 1;
  border-radius: 8px;
  border: 1px solid #e8e8ef;
  object-fit: contain;
  background: #f8f9fb;
  margin: 0 auto;
  display: block;
}

/* Phone popup specific */
.mobile-phone-number {
  font-size: 26px;
  font-weight: 700;
  color: #15294D;
  margin: 12px 0 4px;
  letter-spacing: 1px;
  font-variant-numeric: tabular-nums;
}

.mobile-phone-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.mobile-phone-actions > .public-button {
  flex: 1;
  min-width: 0;
}

.mobile-phone-actions .copied {
  background-color: var(--color-accent-soft);
  transition: none;
}

@media (prefers-reduced-motion: reduce) {
  .cs-tooltip,
  .cs-item:hover .cs-tooltip,
  .cs-item:focus-visible .cs-tooltip {
    transition: none;
  }
}
</style>
