<template>
  <div class="homepage">
    <!-- Hero Carousel Section -->
    <section class="hero-section">
      <div class="hero-carousel">
        <div
          v-for="(slide, index) in heroSlides"
          :key="index"
          class="hero-slide"
          :class="{ active: currentSlide === index }"
          :style="{ backgroundImage: slide.image ? `url(${slide.image})` : undefined }"
        >
          <div class="hero-slide-gradient"></div>
          <div class="hero-glow hero-glow--gold"></div>
          <div class="hero-glow hero-glow--blue"></div>
          <div class="hero-glow hero-glow--amber"></div>
          <div class="hero-glow hero-glow--deep-blue"></div>
          <div class="hero-content container">
            <div class="hero-badge">
              <span>精 品 投 资 移 民 决 策 平 台</span>
            </div>
            <h1 class="hero-title">{{ slide.title }}</h1>
            <p class="hero-subtitle">{{ slide.subtitle }}</p>
            <div class="hero-actions">
              <NuxtLink v-if="slide.link" :to="slide.link" class="btn-hero-primary">
                查看项目详情
              </NuxtLink>
              <NuxtLink to="/contact" class="btn-hero-secondary">
                免费咨询
              </NuxtLink>
            </div>
          </div>
        </div>

        <button class="carousel-arrow carousel-prev" @click="prevSlide" aria-label="Previous slide">
          <span v-html="getIconSvg('chevron-left', 20, 'rgba(255,255,255,0.8)')"></span>
        </button>
        <button class="carousel-arrow carousel-next" @click="nextSlide" aria-label="Next slide">
          <span v-html="getIconSvg('chevron-right', 20, 'rgba(255,255,255,0.8)')"></span>
        </button>

        <div class="carousel-dots">
          <button
            v-for="(slide, index) in heroSlides"
            :key="index"
            class="carousel-dot"
            :class="{ active: currentSlide === index }"
            @click="goToSlide(index)"
            :aria-label="'Go to slide ' + (index + 1)"
          ></button>
        </div>
      </div>
    </section>

    <!-- Trust Bar Section -->
    <section v-if="trustItems.length > 0" ref="trustBarRef" class="trust-bar-section">
      <div class="trust-bar">
        <template v-for="(item, index) in trustItems" :key="index">
          <div class="trust-bar-item">
            <span class="trust-bar-number" :data-target="parseTrustNumber(item.number)">{{ animatedNumbers[index] ?? formatTrustNumber(item.number) }}</span>
            <span class="trust-bar-label">{{ item.label }}</span>
          </div>
          <div v-if="index < trustItems.length - 1" class="trust-bar-divider"></div>
        </template>
      </div>
    </section>

    <!-- Project Cards Section -->
    <section class="section projects-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ sectionTitle }}</h2>
          <p>{{ sectionSubtitle }}</p>
        </div>

        <div v-if="pending.projects" class="loading-state">
          <div class="skeleton" style="height:360px;"></div>
        </div>
        <div v-else-if="error.projects" class="error-state">
          <div class="error-card">
            <span v-html="getIconSvg('alert-circle', 24, '#c0392b')"></span>
            <p>加载失败，请刷新重试</p>
          </div>
        </div>
        <div v-else class="project-cards">
          <div v-for="(project, idx) in projectCards" :key="project.slug" class="project-card reveal">
            <div class="card-image" :class="`card-image--${idx % 3}`">
              <div class="card-image-glow"></div>
              <div class="card-image-overlay"></div>
              <img v-if="project.image" :src="project.image" :alt="project.title" loading="lazy" />
            </div>
            <div class="card-body">
              <h3 class="card-title">{{ project.title }}</h3>
              <p class="card-desc">{{ project.description }}</p>
              <div class="card-stats">
                <span class="card-stat" v-for="(feat, fi) in project.features" :key="fi">
                  {{ feat }}
                </span>
              </div>
              <NuxtLink :to="project.link" class="card-link">
                了解详情
                <span class="link-arrow" v-html="getIconSvg('chevron-right', 14, 'currentColor')"></span>
              </NuxtLink>
            </div>
            <div class="card-bottom-line"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Cases Section -->
    <section class="section cases-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ caseTitle }}</h2>
          <p v-if="caseSubtitle">{{ caseSubtitle }}</p>
        </div>

        <div v-if="pending.cases" class="loading-state">
          <div class="skeleton" style="height:360px;"></div>
        </div>
        <div v-else-if="error.cases" class="error-state">
          <div class="error-card">
            <span v-html="getIconSvg('alert-circle', 24, '#c0392b')"></span>
            <p>加载失败，请刷新重试</p>
          </div>
        </div>
        <div v-else class="cases-grid">
          <CaseCard
            v-for="item in featuredCases"
            :key="item.id"
            :slug="item.slug"
            :name="item.name"
            :country="item.country_from"
            :project="item.project?.name"
            :summary="stripHtml(item.content)"
            :image="item.photo_url"
          />
        </div>

        <div v-if="!pending.cases && featuredCases.length === 0" class="empty-state">
          暂无成功案例
        </div>
      </div>
    </section>

    <!-- Testimonials Section -->
    <section v-if="featuredTestimonials.length > 0" class="section testimonial-section">
      <div class="container">
        <div class="section-header">
          <h2 class="decorate-on">{{ testimonialTitle }}<i class="decorate"></i></h2>
          <p v-if="testimonialSubtitle">{{ testimonialSubtitle }}</p>
        </div>

        <TestimonialCarousel :testimonials="featuredTestimonials" />
      </div>
    </section>

    <!-- Lawyer Team Section -->
    <section v-if="lawyers.length > 0" class="section lawyer-section">
      <div class="container">
        <div class="section-header">
          <h2>专业律师团队</h2>
          <p>资深移民律师，为您保驾护航</p>
        </div>
        <LawyerCarousel :lawyers="lawyers" />
      </div>
    </section>

    <!-- Advantages Section -->
    <section class="section advantages-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ advantageTitle }}</h2>
          <p>{{ advantageSubtitle }}</p>
        </div>

        <div class="advantages-grid">
          <div v-for="(adv, index) in advantages" :key="index" class="advantage-card reveal">
            <div class="advantage-icon">
              <span v-if="getIconByName(adv.icon)" v-html="getIconSvg(adv.icon, 22, '#C8963E')" class="advantage-svg"></span>
              <span v-else class="advantage-svg-fallback">
                <span v-html="getIconSvg('star', 22, '#C8963E')"></span>
              </span>
            </div>
            <h3 class="advantage-title">{{ adv.title }}</h3>
            <p class="advantage-desc">{{ adv.description }}</p>
          </div>
        </div>

        <div v-if="advantageSection?.image" class="advantage-banner">
          <img :src="advantageSection.image" alt="优势区域图" loading="lazy" />
        </div>
      </div>
    </section>

    <!-- CTA Banner -->
    <section class="cta-section">
      <div class="cta-glow cta-glow--gold"></div>
      <div class="cta-glow cta-glow--blue"></div>
      <div class="container">
        <div class="cta-banner">
          <p class="cta-label">Start Your Journey</p>
          <h2 class="cta-title">开启您的移民之旅</h2>
          <p class="cta-desc">专业顾问一对一咨询，为您定制最佳移民方案。首次咨询不收取任何费用。</p>
          <NuxtLink to="/contact" class="btn-primary cta-button">免费咨询顾问</NuxtLink>
          <p class="cta-phone">或拨打 {{ siteConfig?.contact_phone || '400-xxx-xxxx' }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">

useSeo({ title: '首页' });

const { siteConfig } = useSiteConfig();

import { getIconByName, getIconSvg } from '~/composables/lucideIcons'


// Hero carousel
const currentSlide = ref(0);
let autoTimer: ReturnType<typeof setInterval> | null = null;

interface HeroSlide {
  title: string;
  subtitle: string;
  image: string;
  link?: string;
  gradient?: string;
  project_slug?: string;
}

const defaultSlides: HeroSlide[] = [
  {
    title: '美国EB-5投资移民',
    subtitle: '投资80万美元，全家获得美国绿卡',
    image: '',
  },
  {
    title: '香港资本投资者入境计划',
    subtitle: '投资3000万港元，畅享亚洲金融中心',
    image: '',
  },
  {
    title: '巴拿马购房移民',
    subtitle: '30万美元购房，快速获得永久居留权',
    image: '',
  },
];

const heroSlides = ref<HeroSlide[]>(defaultSlides);

// Fetch home config and projects
const { data: homeConfig, pending: pendingHome } = await useFetch('/api/v1/home-config', {
  onResponseError() {
    // Use defaults if API fails
  },
});

const { data: projectsData, pending: pendingProjects, error: errorProjectsRaw } = await useFetch<{
  data?: Array<{ name: string; slug: string; tagline: string; overview_text: string; cover_image: string }>;
}>('/api/v1/projects', {
  onResponseError() {
    // Use defaults if API fails
  },
});

interface CaseItem {
  id: number;
  slug: string;
  name: string;
  country_from: string;
  photo_url: string;
  content: string;
  project?: { name: string };
}

const { data: casesData, pending: pendingCases, error: errorCasesRaw } = await useFetch<{
  data?: CaseItem[];
}>('/api/v1/cases', {
  onResponseError() {},
});

interface LawyerItem {
  id: number;
  photo_url: string;
  name: string;
  title: string;
  tags: string[];
}

interface TestimonialItem {
  id: number;
  nickname: string;
  avatar_url: string;
  rating: number;
  content: string;
}

const { data: testimonialsData } = await useFetch<{ data?: TestimonialItem[] }>('/api/v1/testimonials', {
  onResponseError() {},
});

const { data: lawyersData } = await useFetch<{ data?: LawyerItem[] }>('/api/v1/lawyers', {
  onResponseError() {},
});

const lawyers = computed<LawyerItem[]>(() => {
  const raw = lawyersData.value as { data?: LawyerItem[] } | null;
  return raw?.data ?? [];
});

const pending = computed(() => ({
  projects: pendingProjects.value,
  cases: pendingCases.value,
}));

const error = computed(() => ({
  projects: errorProjectsRaw.value ? '加载失败，请刷新重试' : null,
  cases: errorCasesRaw.value ? '加载失败，请刷新重试' : null,
}));

// Override slides from API if available
if (homeConfig.value) {
  const config = homeConfig.value as unknown as Record<string, unknown>;
  const data = config.data as Record<string, unknown> | undefined;
  if (data && Array.isArray(data.hero_slides)) {
    heroSlides.value = (data.hero_slides as Array<Record<string, string>>).map((s) => ({
      title: s.title || '',
      subtitle: s.desc || '',
      image: s.image || '',
      link: s.link || (s.project_slug ? `/projects/${s.project_slug}` : ''),
      gradient: s.gradient || '',
      project_slug: s.project_slug || '',
    }));
  }
}

const showcaseConfig = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.project_showcase) {
      return data.project_showcase as {
        section_title?: string;
        section_subtitle?: string;
        featured_slugs?: string[];
      };
    }
  }
  return null;
});

const caseShowcaseConfig = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.case_showcase) {
      return data.case_showcase as {
        section_title?: string;
        section_subtitle?: string;
        featured_case_ids?: number[];
      };
    }
  }
  return null;
});

const caseTitle = computed(() => caseShowcaseConfig.value?.section_title || '成功案例');
const caseSubtitle = computed(() => caseShowcaseConfig.value?.section_subtitle || '');

function stripHtml(html: string): string {
  if (!html) return '';
  return html.replace(/<[^>]+>/g, '').replace(/&nbsp;/g, ' ').slice(0, 80);
}

const featuredCases = computed<CaseItem[]>(() => {
  const apiData = casesData.value as { data?: CaseItem[] } | null;
  const all = apiData?.data ?? [];
  if (all.length === 0) return [];

  const featured = caseShowcaseConfig.value?.featured_case_ids;
  if (featured && featured.length > 0) {
    const orderMap = new Map(featured.map((id: number, i: number) => [id, i]));
    return all
      .filter((c) => orderMap.has(c.id))
      .sort((a, b) => {
        const ai = orderMap.get(a.id);
        const bi = orderMap.get(b.id);
        if (ai !== undefined && bi !== undefined) return ai - bi;
        return 0;
      });
  }

  return all;
});

const testimonialShowcaseConfig = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.testimonial_showcase) {
      return data.testimonial_showcase as {
        section_title?: string;
        section_subtitle?: string;
        featured_testimonial_ids?: number[];
      };
    }
  }
  return null;
});

const testimonialTitle = computed(() => testimonialShowcaseConfig.value?.section_title || '客户评价');
const testimonialSubtitle = computed(() => testimonialShowcaseConfig.value?.section_subtitle || '');

const featuredTestimonials = computed<TestimonialItem[]>(() => {
  const apiData = testimonialsData.value as { data?: TestimonialItem[] } | null;
  const all = apiData?.data ?? [];
  if (all.length === 0) return [];

  const featured = testimonialShowcaseConfig.value?.featured_testimonial_ids;
  if (featured && featured.length > 0) {
    const orderMap = new Map(featured.map((id: number, i: number) => [id, i]));
    return all
      .filter((t) => orderMap.has(t.id))
      .sort((a, b) => {
        const ai = orderMap.get(a.id);
        const bi = orderMap.get(b.id);
        if (ai !== undefined && bi !== undefined) return ai - bi;
        return 0;
      });
  }

  return all;
});

const advantageSection = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.advantage_section) {
      return data.advantage_section as { section_title?: string; section_subtitle?: string; image?: string };
    }
  }
  return null;
});

const trustItems = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && Array.isArray(data.hero_trust)) {
      return data.hero_trust as Array<{ number: string; label: string }>;
    }
  }
  return [];
});

const trustBarRef = ref<HTMLElement | null>(null);
const animatedNumbers = ref<string[]>([]);
let trustAnimating = false;

function parseTrustNumber(raw: string): number {
  const match = raw.replace(/,/g, '').match(/([\d.]+)/);
  return match ? parseInt(match[1], 10) : 0;
}

function formatTrustNumber(raw: string): string {
  return raw;
}

function startTrustCountUp(items: Array<{ number: string }>) {
  if (trustAnimating) return;
  trustAnimating = true;

  const targets = items.map((item) => parseTrustNumber(item.number));
  const suffixes = items.map((item) => {
    const num = parseTrustNumber(item.number);
    const raw = item.number.replace(/,/g, '');
    return raw.replace(String(num), '').replace(/^\d+/, '');
  });

  const duration = 800;
  const start = performance.now();

  function step(timestamp: number) {
    const elapsed = timestamp - start;
    const progress = Math.min(elapsed / duration, 1);
    const eased = progress === 1 ? 1 : 1 - Math.pow(2, -10 * progress);

    animatedNumbers.value = targets.map((t, i) => {
      const val = Math.floor(eased * t);
      return val.toLocaleString() + suffixes[i];
    });

    if (progress < 1) {
      requestAnimationFrame(step);
    }
  }

  requestAnimationFrame(step);
}

const advantageTitle = computed(() => advantageSection.value?.section_title || '为什么选择北极星移民');
const advantageSubtitle = computed(() => advantageSection.value?.section_subtitle || '专业服务，值得信赖');

const sectionTitle = computed(() => showcaseConfig.value?.section_title || '精选移民项目');
const sectionSubtitle = computed(() => showcaseConfig.value?.section_subtitle || '为您量身定制的最佳移民方案');

interface ProjectCard {
  slug: string;
  title: string;
  description: string;
  image: string;
  features: string[];
  link: string;
}

const projectCards = computed<ProjectCard[]>(() => {
  const apiProjects = (projectsData.value as unknown as {
    data?: Array<{ name: string; slug: string; tagline: string; overview_text: string; cover_image: string }>;
  })?.data;

  if (apiProjects && apiProjects.length > 0) {
    const featured = showcaseConfig.value?.featured_slugs;
    let items: ProjectCard[] = apiProjects.map((p) => ({
      slug: p.slug,
      title: p.name,
      description: p.tagline || p.overview_text || '',
      image: p.cover_image || '',
      features: [],
      link: `/projects/${p.slug}`,
    }));

    if (featured && featured.length > 0) {
      const orderMap = new Map(featured.map((s: string, i: number) => [s, i]));
      items.sort((a, b) => {
        const ai = orderMap.get(a.slug);
        const bi = orderMap.get(b.slug);
        if (ai !== undefined && bi !== undefined) return ai - bi;
        if (ai !== undefined) return -1;
        if (bi !== undefined) return 1;
        return 0;
      });
    }

    return items;
  }

  return [
    {
      slug: 'eb5',
      title: '美国EB-5投资移民',
      description: '投资80万美元到美国政府批准的商业项目，创造10个就业机会，全家获得美国绿卡。',
      image: '',
      features: ['投资金额：80万美元起', '办理周期：约24-36个月', '适合人群：高净值投资者'],
      link: '/projects/eb5',
    },
    {
      slug: 'cies',
      title: '香港资本投资者入境计划',
      description: '投资3000万港元到香港获许金融资产，即可获得香港居留权，7年后可申请永久居民。',
      image: '',
      features: ['投资金额：3000万港元', '办理周期：约6-12个月', '适合人群：资产雄厚人士'],
      link: '/projects/cies',
    },
    {
      slug: 'panama',
      title: '巴拿马购房移民',
      description: '购买30万美元以上巴拿马房产，快速获得巴拿马永久居留权，享受中美洲优质生活。',
      image: '',
      features: ['投资金额：30万美元起', '办理周期：约3-6个月', '适合人群：有海外置业需求者'],
      link: '/projects/panama',
    },
  ];
});

const advantages = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && Array.isArray(data.advantage_items)) {
      return (data.advantage_items as Array<Record<string, string>>).map((a) => ({
        icon: a.icon || '',
        iconType: a.icon_type || '',
        title: a.title || '',
        description: a.description || '',
      }));
    }
  }

  return [
    { icon: '\u{1F3C6}', iconType: '', title: '10年行业经验', description: '深耕投资移民领域，拥有丰富的成功案例和行业资源' },
    { icon: '\u{1F465}', iconType: '', title: '专业顾问团队', description: '资深移民律师、前移民官组成的一流专业团队' },
    { icon: '\u{1F512}', iconType: '', title: '100%成功率', description: '严格的审核流程确保申请质量，保持行业领先成功率' },
    { icon: '\u{1F4DE}', iconType: '', title: '一站式服务', description: '从方案定制到成功获批，全程跟踪服务让您无忧' },
  ];
});

// Carousel logic
const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % heroSlides.value.length;
  resetAutoPlay();
};

const prevSlide = () => {
  currentSlide.value =
    (currentSlide.value - 1 + heroSlides.value.length) % heroSlides.value.length;
  resetAutoPlay();
};

const goToSlide = (index: number) => {
  currentSlide.value = index;
  resetAutoPlay();
};

const resetAutoPlay = () => {
  if (autoTimer) clearInterval(autoTimer);
  autoTimer = setInterval(() => {
    currentSlide.value = (currentSlide.value + 1) % heroSlides.value.length;
  }, 5000);
};

let revealObserver: IntersectionObserver | null = null;

onMounted(() => {
  autoTimer = setInterval(() => {
    currentSlide.value = (currentSlide.value + 1) % heroSlides.value.length;
  }, 5000);

  revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          revealObserver?.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.15, rootMargin: '0px 0px -30px 0px' }
  );

  document.querySelectorAll('.reveal').forEach((el) => revealObserver!.observe(el));
});

// Trust bar count-up observer
let trustObserver: IntersectionObserver | null = null;

watch([trustBarRef, trustItems], ([el, items]) => {
  if (el && items.length > 0 && !trustAnimating) {
    trustObserver?.disconnect();
    trustObserver = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          startTrustCountUp(items);
          trustObserver?.disconnect();
        }
      },
      { threshold: 0.3 }
    );
    trustObserver.observe(el as HTMLElement);
  }
});

onUnmounted(() => {
  if (autoTimer) clearInterval(autoTimer);
  if (revealObserver) revealObserver.disconnect();
  if (trustObserver) trustObserver.disconnect();
});
</script>

<style scoped>
/* Hero Section */
.hero-section {
  /* hero starts below fixed header via .main-content margin-top */
}

.hero-carousel {
  position: relative;
  height: 560px;
  overflow: hidden;
  background: linear-gradient(135deg, #080E1A 0%, #0F1E3D 35%, #15294D 65%, #1A3A5C 100%);
}

.hero-slide {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  opacity: 0;
  transition: opacity 0.6s ease-in-out;
}

.hero-slide.active {
  opacity: 1;
}

.hero-slide-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(8, 14, 26, 0.7) 0%, rgba(21, 41, 77, 0.4) 50%, rgba(26, 58, 92, 0.6) 100%);
}

.hero-glow {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  filter: blur(60px);
  will-change: transform, opacity, border-radius;
}

.hero-glow--gold {
  top: -120px;
  right: -80px;
  width: 520px;
  height: 520px;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.15) 0%, rgba(200, 150, 62, 0.04) 30%, transparent 65%);
  animation: blobFloat1 12s ease-in-out infinite alternate,
             blobPulse 8s ease-in-out infinite,
             blobMorph 14s ease-in-out infinite;
  animation-delay: 0s, 0s, 0s;
}

.hero-glow--blue {
  bottom: -100px;
  left: -60px;
  width: 380px;
  height: 380px;
  background: radial-gradient(circle, rgba(30, 80, 160, 0.5) 0%, rgba(30, 58, 110, 0.15) 40%, transparent 70%);
  animation: blobFloat2 14s ease-in-out infinite alternate,
             blobPulse 10s ease-in-out infinite 2s,
             blobMorph 16s ease-in-out infinite 3s;
}

.hero-glow--amber {
  top: 50%;
  right: -40px;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(220, 170, 70, 0.18) 0%, rgba(200, 140, 50, 0.05) 30%, transparent 65%);
  animation: blobFloat3 16s ease-in-out infinite alternate,
             blobPulse 9s ease-in-out infinite 4s,
             blobMorph 18s ease-in-out infinite 1s;
}

.hero-glow--deep-blue {
  top: -40px;
  left: 20%;
  width: 240px;
  height: 240px;
  background: radial-gradient(circle, rgba(20, 40, 80, 0.6) 0%, rgba(15, 30, 60, 0.2) 35%, transparent 65%);
  animation: blobFloat4 15s ease-in-out infinite alternate,
             blobPulse 11s ease-in-out infinite 1s,
             blobMorph 13s ease-in-out infinite 5s;
}

/* ── Ambient Light Blob Keyframes ── */

@keyframes blobFloat1 {
  0%   { transform: translate(0, 0) rotate(0deg); }
  33%  { transform: translate(30px, -20px) rotate(120deg); }
  66%  { transform: translate(-15px, 25px) rotate(240deg); }
  100% { transform: translate(20px, -10px) rotate(360deg); }
}

@keyframes blobFloat2 {
  0%   { transform: translate(0, 0) rotate(0deg); }
  33%  { transform: translate(-25px, -15px) rotate(120deg); }
  66%  { transform: translate(20px, 20px) rotate(240deg); }
  100% { transform: translate(-10px, -25px) rotate(360deg); }
}

@keyframes blobFloat3 {
  0%   { transform: translate(0, 0) rotate(0deg); }
  33%  { transform: translate(-20px, 25px) rotate(140deg); }
  66%  { transform: translate(25px, -15px) rotate(260deg); }
  100% { transform: translate(-15px, 10px) rotate(380deg); }
}

@keyframes blobFloat4 {
  0%   { transform: translate(0, 0) rotate(0deg); }
  33%  { transform: translate(15px, -30px) rotate(100deg); }
  66%  { transform: translate(-25px, -10px) rotate(220deg); }
  100% { transform: translate(10px, 20px) rotate(340deg); }
}

@keyframes blobPulse {
  0%, 100% { opacity: 0.5; }
  50%      { opacity: 0.9; }
}

@keyframes blobMorph {
  0%   { border-radius: 50% 50% 50% 50%; }
  25%  { border-radius: 45% 55% 48% 52%; }
  50%  { border-radius: 55% 45% 52% 48%; }
  75%  { border-radius: 48% 52% 55% 45%; }
  100% { border-radius: 50% 50% 50% 50%; }
}

.hero-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  color: var(--bg-white);
}

.hero-badge {
  margin-bottom: 20px;
}

.hero-badge span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 14px;
  background: rgba(200, 150, 62, 0.08);
  border: 1px solid rgba(200, 150, 62, 0.15);
  border-radius: 20px;
  color: var(--accent);
  font-size: 12px;
  letter-spacing: 2px;
}

.hero-title {
  font-family: var(--font-serif);
  font-size: 44px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 16px;
  max-width: 520px;
}

.hero-subtitle {
  font-size: 16px;
  opacity: 0.65;
  line-height: 1.7;
  margin-bottom: 32px;
  max-width: 420px;
}

.hero-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 40px;
}

.btn-hero-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  padding: 13px 32px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  color: var(--bg-white);
  border: none;
  border-radius: var(--radius-md);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: var(--shadow-gold);
  transition: all 0.2s ease;
  overflow: hidden;
}

.btn-hero-primary::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 60%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
  transform: skewX(-25deg);
  animation: shimmer 3s ease-in-out infinite 1s;
}

.btn-hero-primary:hover {
  box-shadow: 0 6px 24px rgba(200, 150, 62, 0.4);
  transform: translateY(-1px);
}

.btn-hero-primary:hover::after {
  animation-duration: 1.5s;
}

@keyframes shimmer {
  0%   { left: -100%; }
  100% { left: 150%; }
}

.btn-hero-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 13px 28px;
  background: rgba(255, 255, 255, 0.06);
  color: var(--bg-white);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: var(--radius-md);
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-hero-secondary:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.35);
}

/* Trust Bar Section */
.trust-bar-section {
  background: rgba(15, 27, 45, 0.82);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.trust-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 64px;
  padding: 36px 40px;
  max-width: 960px;
  margin: 0 auto;
}

.trust-bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.trust-bar-number {
  font-family: var(--font-serif);
  font-size: 36px;
  font-weight: 700;
  color: var(--color-accent);
  letter-spacing: -0.5px;
  font-variant-numeric: tabular-nums;
}

.trust-bar-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 1px;
}

.trust-bar-divider {
  width: 1px;
  height: 40px;
  background: linear-gradient(180deg, transparent, rgba(200, 150, 62, 0.3), transparent);
}

@media (max-width: 767px) {
  .trust-bar {
    flex-direction: column;
    gap: 20px;
    padding: 28px 20px;
  }

  .trust-bar-divider {
    width: 60px;
    height: 1px;
  }

  .trust-bar-number {
    font-size: 30px;
  }
}

.carousel-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  backdrop-filter: blur(4px);
  transition: all 0.2s ease;
  z-index: 2;
  padding: 0;
}

.carousel-arrow:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.25);
}

.carousel-prev {
  left: 24px;
}

.carousel-next {
  right: 24px;
}

.carousel-dots {
  position: absolute;
  bottom: 28px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 20px;
  z-index: 2;
}

.carousel-dot {
  width: 24px;
  height: 2px;
  border-radius: 1px;
  border: none;
  background: rgba(255, 255, 255, 0.2);
  cursor: pointer;
  padding: 0;
  transition: all 0.25s ease-out;
}

.carousel-dot.active {
  background: var(--accent);
  width: 32px;
}

@media (max-width: 1023px) {
  .hero-carousel {
    height: 480px;
  }

  .hero-title {
    font-size: 36px;
    max-width: 440px;
  }
}

@media (max-width: 767px) {
  .hero-carousel {
    height: 460px;
  }

  .hero-title {
    font-size: 32px;
    max-width: 100%;
  }

  .hero-subtitle {
    font-size: 15px;
    max-width: 100%;
  }

  .hero-actions {
    flex-direction: column;
    margin-bottom: 32px;
  }

  .btn-hero-primary,
  .btn-hero-secondary {
    width: 100%;
    justify-content: center;
  }

  .carousel-dot {
    width: 16px;
  }

  .carousel-dot.active {
    width: 24px;
  }
}

/* Sections */
.section {
  padding: var(--section-gap) 0;
}

.section-header {
  text-align: center;
  margin-bottom: 48px;
}

.section-header h2 {
  font-family: var(--font-serif);
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.section-header p {
  font-size: 16px;
  color: var(--text-light);
  margin-bottom: 0;
}

.section-header::after {
  content: '';
  display: block;
  width: 48px;
  height: 2px;
  margin: 16px auto 0;
  background: linear-gradient(90deg, transparent, var(--accent), transparent);
}

/* Project Cards */
.project-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.project-card {
  position: relative;
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--duration-slow) var(--ease-out),
              transform 0.35s var(--ease-spring),
              border-color var(--duration-normal) var(--ease-out);
}

.project-card:hover {
  box-shadow: var(--shadow-xl), 0 0 0 1px rgba(200, 150, 62, 0.25);
  transform: translateY(-6px);
  border-color: rgba(200, 150, 62, 0.4);
}

.card-image {
  height: 200px;
  overflow: hidden;
  position: relative;
}

.card-image-overlay {
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(180deg, transparent 50%, rgba(200, 150, 62, 0.25) 100%);
  opacity: 0;
  transition: opacity var(--duration-slow) var(--ease-out);
}

.card-image::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(to top, rgba(200, 150, 62, 0.3), transparent);
  z-index: 2;
  transition: height var(--duration-slow) var(--ease-out),
              opacity var(--duration-slow) var(--ease-out);
}

.project-card:hover .card-image-overlay {
  opacity: 1;
}

.project-card:hover .card-image::after {
  height: 100px;
  background: linear-gradient(to top, rgba(200, 150, 62, 0.4), transparent 80%);
}

.card-image--0 {
  background: linear-gradient(135deg, #0F1E3D, #15294D);
}

.card-image--1 {
  background: linear-gradient(135deg, #15294D, #1A3A5C);
}

.card-image--2 {
  background: linear-gradient(135deg, #1A3A5C, #1E3A6E);
}

.card-image-glow {
  position: absolute;
  top: -30px;
  right: -20px;
  width: 140px;
  height: 140px;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.12), transparent 70%);
  border-radius: 50%;
  z-index: 1;
  transition: all var(--duration-slow) var(--ease-out);
}

.project-card:hover .card-image-glow {
  width: 220px;
  height: 220px;
  top: -60px;
  right: -60px;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.22), transparent 65%);
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.45s var(--ease-out);
}

.project-card:hover .card-image img {
  transform: scale(1.08);
}

.card-body {
  padding: 22px 24px;
  position: relative;
  z-index: 2;
}

.card-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8px;
  transition: color var(--duration-normal) var(--ease-out);
}

.project-card:hover .card-title {
  color: var(--color-accent-dark);
}

.card-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.7;
  margin-bottom: 16px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-stats {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 18px;
}

.card-stat {
  background: var(--color-accent-soft);
  color: var(--color-accent-dark);
  font-size: 12px;
  font-weight: 500;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  transition: background var(--duration-normal) var(--ease-out),
              color var(--duration-normal) var(--ease-out);
}

.project-card:hover .card-stat {
  background: rgba(200, 150, 62, 0.18);
  color: var(--color-accent-dark);
}

.card-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  font-weight: 600;
  color: var(--primary);
  transition: gap var(--duration-normal) var(--ease-out),
              color var(--duration-normal) var(--ease-out);
}

.card-link:hover {
  gap: 8px;
  color: var(--accent-dark);
}

.link-arrow {
  display: inline-flex;
  align-items: center;
  transition: transform var(--duration-normal) var(--ease-out);
}

.card-link:hover .link-arrow {
  transform: translateX(2px);
}

/* ── Bottom golden line ── */

.card-bottom-line {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--gradient-gold);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.45s var(--ease-spring);
  z-index: 3;
}

.project-card:hover .card-bottom-line {
  transform: scaleX(1);
}

@media (max-width: 1023px) {
  .project-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .project-cards {
    grid-template-columns: 1fr;
  }
}

/* Advantages */
.advantages-section {
  background-color: var(--bg-light);
}

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.advantage-card {
  text-align: center;
  padding: 28px 20px;
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  position: relative;
  overflow: hidden;
  transition: box-shadow 0.3s ease, border-color 0.3s ease;
}

.advantage-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--accent), transparent);
}

.advantage-card:hover {
  box-shadow: var(--shadow-md);
  border-color: rgba(200, 150, 62, 0.15);
}

.advantage-icon {
  width: 52px;
  height: 52px;
  margin: 0 auto 14px;
  background: linear-gradient(135deg, rgba(21, 41, 77, 0.06), rgba(21, 41, 77, 0.02));
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.advantage-svg {
  display: flex;
  align-items: center;
  justify-content: center;
}

.advantage-svg-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
}

.advantage-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.advantage-desc {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.advantage-banner {
  margin-top: 32px;
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.advantage-banner img {
  width: 100%;
  display: block;
  border-radius: var(--radius-lg);
}

@media (max-width: 1023px) {
  .advantages-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .advantages-grid {
    grid-template-columns: 1fr;
  }
}

/* Cases Section */
.cases-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-light);
  font-size: 15px;
}

/* Testimonial Section */
/* Testimonial Carousel */
.testimonial-section {
  background: var(--bg-white);
}

.section-header .decorate-on {
  position: relative;
  display: inline-block;
}

.section-header .decorate {
  display: block;
  width: 48px;
  height: 3px;
  background: var(--primary);
  margin: 10px auto 0;
  border-radius: 2px;
}

/* CTA Section */
.cta-section {
  padding: 72px 0;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #0F1E3D 0%, #15294D 40%, #1A3A5C 100%);
}

.cta-glow {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.cta-glow--gold {
  top: 0;
  left: 30%;
  width: 500px;
  height: 500px;
  background: radial-gradient(ellipse, rgba(200, 150, 62, 0.06) 0%, transparent 60%);
  transform: translateY(-50%);
}

.cta-glow--blue {
  bottom: 0;
  right: 20%;
  width: 400px;
  height: 400px;
  background: radial-gradient(ellipse, rgba(30, 58, 110, 0.5) 0%, transparent 50%);
  transform: translateY(50%);
}

.cta-banner {
  text-align: center;
  color: var(--bg-white);
  position: relative;
  z-index: 1;
}

.cta-label {
  font-size: 13px;
  color: var(--accent);
  letter-spacing: 3px;
  text-transform: uppercase;
  margin-bottom: 12px;
}

.cta-title {
  font-family: var(--font-serif);
  font-size: 34px;
  font-weight: 700;
  margin-bottom: 12px;
}

.cta-desc {
  font-size: 15px;
  opacity: 0.65;
  margin-bottom: 32px;
  max-width: 420px;
  margin-left: auto;
  margin-right: auto;
}

.cta-button {
  font-size: 16px;
  padding: 14px 36px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-gold);
  transition: all 0.2s ease;
}

.cta-button:hover {
  box-shadow: 0 6px 24px rgba(200, 150, 62, 0.4);
  transform: translateY(-1px);
}

.cta-phone {
  font-size: 13px;
  opacity: 0.45;
  margin-top: 14px;
}

.loading-state {
  text-align: center;
  padding: 40px;
}

.error-state {
  text-align: center;
  padding: 40px;
}

.error-card {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 28px 40px;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  color: var(--text-light);
  font-size: 15px;
}

.error-card span {
  display: flex;
}

/* Responsive */
@media (max-width: 1023px) {
  .section {
    padding: var(--section-gap-mobile) 0;
  }

  .section-header {
    margin-bottom: 32px;
  }

  .section-header h2 {
    font-size: 28px;
  }

  .project-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .cases-grid {
    grid-template-columns: repeat(2, 1fr);
  }



  .advantages-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .section {
    padding: var(--section-gap-mobile) 0;
  }

  .section-header h2 {
    font-size: 24px;
  }

  .hero-carousel {
    height: 400px;
  }

  .hero-title {
    font-size: 32px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .project-cards,
  .cases-grid,
  .advantages-grid {
    grid-template-columns: 1fr;
  }


}

</style>
