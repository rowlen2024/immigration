<template>
  <div class="about-page">
    <div class="container">
      <ProjectBreadcrumb label="公司简介"/>
    </div>

    <!-- ── Hero ── -->
    <section class="hero">
      <div class="container hero-inner">
        <span class="hero-badge">SINCE 2013</span>
        <h1>领航全球<span class="dot">·</span>智启未来</h1>
        <p class="hero-sub">北极星移民，您值得信赖的国际化伙伴</p>
        <p class="hero-desc">深耕移民行业十余载，以专业洞见与诚信服务，为数千家庭精准指引全球身份规划方向</p>
        <NuxtLink to="/contact" class="btn-primary hero-cta">免费咨询</NuxtLink>
      </div>
    </section>

    <!-- ── Stats ── -->
    <section class="stats-section">
      <div class="container">
        <dl class="stats-bar">
          <div class="stat-item reveal"><dt class="stat-lbl">行业经验</dt><dd class="stat-num">15<span class="stat-unit">年+</span></dd></div>
          <div class="stat-item reveal d1"><dt class="stat-lbl">成功案例</dt><dd class="stat-num">20,000<span class="stat-unit">+</span></dd></div>
          <div class="stat-item reveal d2"><dt class="stat-lbl">申请成功率</dt><dd class="stat-num">96.67<span class="stat-unit">%</span></dd></div>
          <div class="stat-item reveal d2"><dt class="stat-lbl">考察国家和地区</dt><dd class="stat-num">60<span class="stat-unit">+</span></dd></div>
          <div class="stat-item reveal d3"><dt class="stat-lbl">EB5 获批案例</dt><dd class="stat-num">290<span class="stat-unit">+</span></dd></div>
        </dl>
      </div>
    </section>

    <!-- ── Intro ── -->
    <section class="section">
      <div class="container">
        <div class="intro-grid">
          <div class="intro-text reveal">
            <span class="section-label">关于我们</span>
            <p>北极星移民是深圳福瑞特出入境信息咨询有限公司运营的移民服务品牌，成立于<strong>2013年</strong>，总部位于深圳，专注于海外身份规划与全球资产配置等高端移民服务。</p>
            <p>作为移民领域的权威机构，北极星移民汇聚了一支由资深移民顾问、专业项目风控团队、文案律师团队、税务专家及跨国教育顾问组成的精英团队。他们凭借对各国移民政策的深刻解读与丰富的实操经验，为客户精准匹配最优路径。</p>
            <p>作为<strong>美国EB-5投资移民、香港投资移民、巴拿马购房投资移民，美国EB-1A杰出人才移民、还是欧洲黄金签证、加勒比海等各国护照的先驱者</strong>，北极星向来以数据化分析、个性化方案及全程透明化管理，确保成功率与效率的双重保障。同时，我们更提供涵盖海外置业、税务优化、国际教育等增值服务，真正实现"移民+安家+发展"的一站式解决方案。</p>
          </div>
          <div class="intro-img reveal">
            <img :src="`/static/about/images/1778568612114621897.jpg`" alt="北极星移民团队" width="800" height="600" loading="lazy" />
          </div>
        </div>
      </div>
    </section>

    <!-- ── Milestones ── -->
    <section class="section section-alt">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">行业创举</span>
          <h2>行业创举引领，树立专业新标杆</h2>
        </div>

        <div class="card-grid cols-3">
          <div class="card reveal">
            <div class="card-icon" v-html="icons.clock" />
            <h3>15年+行业深耕</h3>
            <p>核心团队成员实地考察各个移民项目，走访60+国家和地区，2万+成功案例。</p>
          </div>
          <div class="card reveal d1">
            <div class="card-icon" v-html="icons.diamond" />
            <h3>巴拿马项目先驱</h3>
            <p>中国大陆首批推广巴拿马购房永居项目的移民机构之一，也是整个行业最早拥有巴拿马项目成功案例的移民公司。</p>
          </div>
          <div class="card reveal d2">
            <div class="card-icon" v-html="icons.shield" />
            <h3>EB5 全流程经验</h3>
            <p>全国少有的EB5全流程经验企业，获批案例290+。美国EB-1A项目大陆地区8天获批纪录保持者，70岁申请人获批年龄极限突破机构。</p>
          </div>
          <div class="card reveal">
            <div class="card-icon" v-html="icons.check" />
            <h3>香港成功率90%</h3>
            <p>香港投资移民、高才通等多项通道全覆盖，续签成功率100%。</p>
          </div>
          <div class="card reveal d1">
            <div class="card-icon" v-html="icons.trending" />
            <h3>成功率 96.67%</h3>
            <p>所有项目申请成功率96.67%，业内高水准。投资类移民成功率100%，身份办理通过率95%+。</p>
          </div>
          <div class="card reveal d2">
            <div class="card-icon" v-html="icons.book" />
            <h3>教育子品牌"三只小象"</h3>
            <p>独家冠名4届大湾区国际教育展，全网粉丝10万+，首创"身份+教育+资产"三维服务体系。</p>
          </div>
        </div>

        <div class="img-pair reveal">
          <figure>
            <img :src="`/static/about/images/1778568651667806803.jpg`" alt="北极星移民成功案例墙" width="800" height="500" loading="lazy" />
            <figcaption>▲北极星移民成功案例墙</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/1778568704715310459.jpg`" alt="巴拿马论坛合影" width="800" height="500" loading="lazy" />
            <figcaption>▲与巴拿马合格投资者法案修订人RICHARD THOMAS在巴拿马合格投资者永居计划巡回论坛合影</figcaption>
          </figure>
        </div>

        <div class="dark-card reveal">
          <h4>战略孵化教育子品牌"三只小象"</h4>
          <ul>
            <li>独家冠名4届大湾区规模最大的国际教育展</li>
            <li>拥有大湾区最具影响力的全媒体平台，全网粉丝10万＋</li>
            <li>全网唯一一家探校广州民心学校平台</li>
            <li>首创"身份+教育+资产"三维服务体系，直击家庭核心诉求，让移民真正成为高净值家庭进阶的跳板</li>
          </ul>
        </div>
      </div>
    </section>

    <!-- ── Honors ── -->
    <section class="section">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">企业荣誉</span>
          <h2>荣誉满载，见证实力与口碑</h2>
          <p>每一个奖状、每一份证书，都是社会以及项目机构对北极星移民坚守专业、诚信服务的最好背书。</p>
        </div>
        <div class="honors-grid reveal">
          <figure v-for="(h, i) in honorImgs" :key="h.src" tabindex="0" role="button" :aria-label="`查看${h.alt}大图`" @click="selectedCert = i" @keydown.enter="selectedCert = i" @keydown.space.prevent="selectedCert = i">
            <img :src="`/static/about/images/${h.src}`" :alt="h.alt" width="600" height="800" loading="lazy" />
          </figure>
        </div>
      </div>
    </section>

    <!-- ── Advantages ── -->
    <section class="section section-alt">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">北极星优势</span>
          <h2>专业、诚信、安心之选</h2>
          <p>选择北极星移民｜意味着选择了一套难以复制的核心优势组合</p>
        </div>
        <div class="adv-grid">
          <div v-for="(a, i) in advantages" :key="a.title" class="adv-card reveal" :class="`d${i % 3}`">
            <div class="adv-icon" v-html="a.icon" />
            <h3>{{ a.title }}</h3>
            <ul>
              <li v-for="p in a.points" :key="p">{{ p }}</li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Founder ── -->
    <section class="section">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">掌舵者之光</span>
          <h2>创始人的远见与实践</h2>
        </div>

        <div class="founder-profile reveal">
          <div class="founder-img">
            <img :src="`/static/about/images/1778576160647356364.jpg`" alt="创始人夏竹青先生" width="600" height="800" loading="lazy" />
          </div>
          <div class="founder-info">
            <h3>夏竹青 先生</h3>
            <p class="founder-title">创始人兼董事长</p>
            <p>拥有<strong>16年国际移民与海外身份规划</strong>的丰富经验，亲身经历行业从萌芽到发展、规范的全过程，对全球移民政策演变、各国项目优劣、市场风险痛点有着深刻独到的理解与前瞻性判断。</p>
            <div class="founder-quote"><p>"坚持每年投入大量时间，亲自带领团队飞赴全球各个项目地实地考察。用脚步丈量项目，只荐真项目、只说真话。"</p></div>
            <p>他将<strong>"诚信专业·使命必达"</strong>定为北极星不可动摇的服务宗旨。要求团队对<strong>每一个案例精雕细琢</strong>，对每一份文件严谨审核，对每一位客户坦诚相待。</p>
          </div>
        </div>

        <div class="founder-gallery reveal">
          <img :src="`/static/about/images/1778576174445265385.jpg`" alt="夏竹青先生活动" width="400" height="300" loading="lazy" />
          <img :src="`/static/about/images/1778576182506577294.jpg`" alt="夏竹青先生考察" width="400" height="300" loading="lazy" />
          <img :src="`/static/about/images/1778576190542151423.jpg`" alt="夏竹青先生演讲" width="400" height="300" loading="lazy" />
          <img :src="`/static/about/images/1778576267975359424.jpg`" alt="夏竹青先生团队" width="400" height="300" loading="lazy" />
        </div>

        <hr class="founder-divider" />

        <div class="founder-profile reveal">
          <div class="founder-img">
            <img :src="`/static/about/images/1778576213314677461.jpg`" alt="联合创始人张霞女士" width="600" height="800" loading="lazy" />
          </div>
          <div class="founder-info">
            <h3>张霞 女士</h3>
            <p class="founder-title">联合创始人 · 运营总监</p>
            <p>拥有<strong>18年移民行业经验及7年国际教育行业经验</strong>，丰富的跨界经历让她对客户在移民后的教育规划等需求有着深刻的理解。</p>
            <p>2020年，她创办了个人品牌<strong>《K妈谈择校》</strong>，为众多移民家庭在子女教育选择上提供了宝贵的指导。她秉承北极星移民的服务宗旨，打造出公司完善的服务体系，使得北极星移民四大部门项目部、顾问部、文案部、客户部协同提供高品质高效率服务。</p>
            <p>正是在夏竹青先生和张霞女士的共同带领下，北极星移民才能够在行业内不断创新，取得诸多成就，积累了服务客户的种种优势，赢得了客户的广泛赞誉和信赖。</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Business ── -->
    <section class="section section-alt">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">优势业务</span>
          <h2>多元化、战略性布局全球，让移民更从容</h2>
        </div>

        <div class="biz-grid">
          <div class="biz-card reveal">
            <div class="biz-num">01</div>
            <h3>巴拿马购房移民</h3>
            <p class="biz-subtitle">项目亮点</p>
            <ul>
              <li>投资30万美元购房，30个工作日速获永久居留卡</li>
              <li>无居住要求，仅需每两年登陆一次</li>
              <li>房产持有满5年后可自由出售，不影响身份持续性</li>
              <li>永居卡无有效期，无需续签</li>
              <li>三代同移，持永居满5年可入籍，免签148国</li>
            </ul>
            <img :src="`/static/about/images/1778576702827841208.png`" alt="巴拿马项目详情" width="640" height="360" loading="lazy" />
            <p class="biz-highlight">国内首个获批案例：上海外贸企业主W先生，通过北极星购入巴拿马城核心区公寓，仅14天成功获批，同步解决税务与子女华侨生联考需求。</p>
          </div>

          <div class="biz-card reveal d1">
            <div class="biz-num">02</div>
            <h3>美国移民</h3>
            <p class="biz-subtitle">人才投资通道全覆盖</p>
            <ul>
              <li><strong>EB-5投资移民：</strong>新政红利期，无排期，双递交，投资80万美元起</li>
              <li><strong>NIW国家利益豁免：</strong>低门槛，重点考察学历或专业能力，前无被拒案例</li>
              <li><strong>EB-1A杰出人才移民：</strong>无需雇主担保，98%客户无补件获批</li>
              <li>合作项目均获移民局预审批（I-924范例认证），资金安全双保险</li>
            </ul>
          </div>

          <div class="biz-card reveal">
            <div class="biz-num">03</div>
            <h3>希腊购房移民</h3>
            <p class="biz-subtitle">欧盟身份快速通道</p>
            <ul>
              <li>购置不动产不低于25万/40万/80万欧元，可获批希腊居留许可证</li>
              <li>持有房产期间可续签5年期签证，房产可出租</li>
              <li>甄选业内资深开发商，品质卓越、合规可信</li>
              <li>专业团队亲赴项目实地考察，严谨筛选每一套房源</li>
              <li>自2013年起深耕希腊项目，历经多轮政策变迁，积累深厚经验</li>
            </ul>
          </div>

          <div class="biz-card reveal d1">
            <div class="biz-num">04</div>
            <h3>中国香港身份</h3>
            <p class="biz-subtitle">多通道 · 高成功率</p>
            <ul>
              <li><strong>投资移民：</strong>投资3000万港币，成熟流程，超高获批成功率</li>
              <li><strong>高才通A类：</strong>年收入250万港币，独创收入合规重组方案</li>
              <li><strong>高才通B/C类：</strong>百强名校毕业生，挖掘隐性贡献与竞争力</li>
              <li><strong>输入内地人才计划：</strong>拥上千优质雇主资源</li>
              <li>香港成功率90%，续签成功率100%</li>
            </ul>
          </div>

          <div class="biz-card reveal">
            <div class="biz-num">05</div>
            <h3>葡萄牙 · 土耳其</h3>
            <p class="biz-subtitle">欧洲身份双通道</p>
            <ul>
              <li><strong>葡萄牙基金移民：</strong>50万欧元，5年后可赎回基金，平均周期降至90天</li>
              <li>合作复星旗下子公司星途ASTER基金，免除中间代理加价</li>
              <li><strong>土耳其房产投资：</strong>40万美元购房，3年后可出售，一步到位获公民身份</li>
              <li><strong>土耳其存款移民：</strong>50万美元定存3年，到期可取回</li>
            </ul>
            <img :src="`/static/about/images/1778576472722322524.jpg`" alt="北极星移民与星途ASTER葡萄牙基金合作授权仪式" width="640" height="360" loading="lazy" />
            <p class="biz-note">▲北极星移民＆星途ASTER葡萄牙基金合作授权仪式</p>
          </div>
        </div>

        <div class="quality-card reveal" style="margin-top:48px">
          <div>
            <h3>项目直控 + 法律双审</h3>
            <p>所有项目由创始人夏竹青带队实地验房；合作律所覆盖美、欧、香港本土牌照律师，确保客户从申请到安家全程无忧。能做到每个项目都能令客户满意，是因为北极星坚持"项目直控+法律双审"模式。</p>
          </div>
          <img :src="`/static/about/images/1778577170456683099.jpg`" alt="北极星移民团队" width="600" height="400" loading="lazy" />
        </div>
      </div>
    </section>

    <!-- ── Closing + CTA ── -->
    <section class="section section-closing">
      <div class="container">
        <div class="closing reveal">
          <h2>以星光为引，共赴璀璨未来</h2>
          <p class="closing-lead">北极星，亘古以来便是天空中的启明星。北极星移民，正是秉承着这份指引的使命，想成为你人生路上的"启明星"。</p>
          <p>我们不能帮助每位客户撰写人生剧本，但我们希望在这个剧本中有那么一个章节是北极星人承载了您的嘱托，帮助您实现新生活，绽放不平凡的生命之光!</p>
          <p>在星辰大海的征途上，北极星移民，永远是您最值得信赖的引航者。</p>
          <p class="closing-cta">让我们携手，以全球视野规划未来，共赴璀璨人生新篇章。</p>
        </div>
      </div>
    </section>

    <ConsultCTA
      title="开启您的全球身份规划"
      description="与资深顾问一对一沟通，获取专属移民方案"
      button-text="免费咨询"
      button-link="/contact"
    />

    <!-- Certificate Lightbox -->
    <Teleport to="body">
      <div v-if="selectedCert !== null" class="cert-lightbox" @click.self="selectedCert = null" @keydown.escape="selectedCert = null">
        <button class="cert-lightbox-close" @click="selectedCert = null" aria-label="关闭">&times;</button>
        <img :src="`/static/about/images/${honorImgs[selectedCert].src}`" :alt="honorImgs[selectedCert].alt" />
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons'

useSeo({
  title: '关于我们',
  description: '北极星移民成立于2013年，专注海外身份规划与全球资产配置。15年+行业经验，2万+成功案例，成功率96.67%。',
  breadcrumbLabel: '关于我们',
})

onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => entries.forEach((e) => { if (e.isIntersecting) e.target.classList.add('visible') }),
    { threshold: 0.15 },
  )
  document.querySelectorAll('.reveal').forEach((el) => observer.observe(el))
})

const honorImgs = [
  { src: '1778637606374142161.png', alt: '北极星移民企业荣誉证书 — 行业资质认证' },
  { src: '1778637708046637702.png', alt: '北极星移民企业荣誉证书 — 优秀服务机构' },
  { src: '1778637719093755072.png', alt: '北极星移民企业荣誉证书 — 诚信经营示范单位' },
  { src: '1778637728035122831.png', alt: '北极星移民企业荣誉证书 — 行业标杆品牌' },
]

const selectedCert = ref<number | null>(null)

const icons: Record<string, string> = {
  clock: getIconSvg('clock', 28),
  diamond: getIconSvg('gem', 28),
  shield: getIconSvg('shield', 28),
  check: getIconSvg('badge-check', 28),
  trending: getIconSvg('trending-up', 28),
  book: getIconSvg('book-open', 28),
}

const advantages = [
  {
    icon: getIconSvg('file-text', 28),
    title: '专业的文案团队',
    points: ['文案团队平均8年以上行业经验', '深度挖掘客户的背景情况', '1000+案子的申请经历'],
  },
  {
    icon: getIconSvg('users', 28),
    title: '行业标杆的顾问团队',
    points: ['团队100%具备海外留学或工作背景，深度了解移民国政策与生活实况', '1V1终身规划，专属顾问全程服务', '精通多国移民项目细则，规避风险盲区'],
  },
  {
    icon: getIconSvg('target', 28),
    title: '风险把控的项目团队',
    points: ['拥有16年行业经验的项目团队', '实地考察60+项目规避投资风险', '至今项目无一列风险存在'],
  },
  {
    icon: getIconSvg('trending-up', 28),
    title: '超高的通过率',
    points: ['全部符合各国移民局精准文件', '申请投资类移民成功率100%', '身份办理通过率95%+'],
  },
  {
    icon: getIconSvg('award', 28),
    title: '优质的后端服务',
    points: ['跟踪移民项目的后端服务', '完善贴心的异国安家服务', '5年+的持续跟进服务'],
  },
]
</script>

<style scoped>
/* ══════ Hero ══════ */
.hero {
  position: relative;
  background: var(--gradient-hero);
  padding: 48px 0 80px;
  overflow: hidden;
  text-align: center;
}
.hero::before, .hero::after {
  content: ''; position: absolute; border-radius: 50%; pointer-events: none;
}
.hero::before {
  top: -180px; right: -160px; width: 520px; height: 520px;
  background: radial-gradient(circle, rgba(200,150,62,.12), transparent 70%);
}
.hero::after {
  bottom: -120px; left: -100px; width: 400px; height: 400px;
  background: radial-gradient(circle, rgba(59,130,246,.08), transparent 70%);
}
.hero-inner { position: relative; z-index: 1; }
.hero-badge {
  display: inline-block; margin-top: 20px; padding: 6px 22px;
  border: 1px solid rgba(200,150,62,.35); border-radius: 999px;
  font-size: 13px; font-weight: 600; letter-spacing: .15em;
  color: var(--accent); background: rgba(200,150,62,.08);
}
.hero h1 {
  font-family: var(--font-serif);
  font-size: clamp(2.4rem, 6vw, 4rem); font-weight: 800;
  color: #fff; margin: 24px 0 16px; line-height: 1.2;
}
.hero h1 .dot { color: var(--accent); }
.hero-sub { font-size: clamp(1.1rem, 2.5vw, 1.35rem); font-weight: 600; color: rgba(255,255,255,.88); margin-bottom: 10px; }
.hero-desc { font-size: 16px; color: rgba(255,255,255,.72); max-width: 600px; margin: 0 auto 28px; line-height: 1.75; }
.hero-cta { display: inline-block; padding: 12px 40px; font-size: 16px; font-weight: 600; border-radius: var(--radius-full); margin-top: 4px; }

/* ══════ Common ══════ */
.section { padding: var(--section-gap) 0; }
.section-alt { background: var(--bg-light); }
.section-label { display: inline-block; font-size: 13px; font-weight: 700; text-transform: uppercase; letter-spacing: .14em; color: var(--accent); margin-bottom: 8px; }
.section-header { text-align: center; margin-bottom: 56px; }
.section-header h2 { font-family: var(--font-serif); font-size: clamp(1.75rem, 4vw, 2.25rem); font-weight: 700; color: var(--text-primary); margin-bottom: 10px; }
.section-header p { font-size: 16px; color: var(--text-light); }
.section-header::after { content: ''; display: block; width: 52px; height: 2px; margin: 18px auto 0; background: linear-gradient(90deg, transparent, var(--accent), transparent); }

/* ══════ Stats ══════ */
.stats-section { background: var(--bg-white); border-bottom: 1px solid var(--border-color); }
.stats-bar { display: flex; justify-content: space-between; padding: 28px 0; gap: 16px; margin: 0; }
.stat-item { text-align: center; flex: 1; }
.stat-num { display: block; font-family: var(--font-serif); font-size: clamp(1.6rem, 3vw, 2.2rem); font-weight: 800; color: var(--primary); line-height: 1.15; margin: 0; }
.stat-unit { font-size: .6em; font-weight: 600; color: var(--accent); }
.stat-lbl { font-size: 13px; color: var(--text-secondary); margin-top: 4px; display: block; font-weight: 400; }

/* ══════ Intro ══════ */
.intro-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 64px; align-items: center; }
.intro-text p { font-size: 16px; color: var(--text-secondary); line-height: 1.9; margin-bottom: 14px; }
.intro-img { border-radius: var(--radius-xl); overflow: hidden; box-shadow: var(--shadow-xl); }
.intro-img img { width: 100%; aspect-ratio: 4/3; object-fit: cover; }

/* ══════ Cards ══════ */
.card-grid { display: grid; gap: 24px; }
.card-grid.cols-3 { grid-template-columns: repeat(3, 1fr); }
.card { background: var(--bg-white); border-radius: var(--radius-lg); padding: 32px 28px; box-shadow: var(--shadow-sm); transition: box-shadow .25s ease, transform .25s ease; }
.card:hover { transform: translateY(-4px); box-shadow: var(--shadow-lg); }
.card-icon { width: 52px; height: 52px; display: flex; align-items: center; justify-content: center; background: var(--accent-soft); border-radius: var(--radius-md); color: var(--accent-dark); margin-bottom: 18px; }
.card h3 { font-family: var(--font-serif); font-size: 18px; font-weight: 700; margin-bottom: 8px; }
.card p { font-size: 14px; color: var(--text-secondary); line-height: 1.7; }

/* ══════ Image pair ══════ */
.img-pair { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin: 48px 0; }
.img-pair figure { border-radius: var(--radius-lg); overflow: hidden; box-shadow: var(--shadow-md); }
.img-pair img { width: 100%; aspect-ratio: 16/10; object-fit: cover; }
.img-pair figcaption { text-align: center; font-size: 13px; color: var(--text-light); padding: 12px; background: var(--bg-white); }

/* ══════ Dark card ══════ */
.dark-card { background: var(--gradient-hero); border-radius: var(--radius-xl); padding: 44px 48px; color: #fff; }
.dark-card h4 { font-family: var(--font-serif); font-size: 1.35rem; font-weight: 700; margin-bottom: 18px; }
.dark-card li { position: relative; padding-left: 22px; margin-bottom: 10px; font-size: 15px; color: rgba(255,255,255,.78); line-height: 1.65; }
.dark-card li::before { content: ''; position: absolute; left: 0; top: 10px; width: 6px; height: 6px; border-radius: 50%; background: var(--accent); }

/* ══════ Advantages ══════ */
.adv-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; }
.adv-card { background: var(--bg-light); border-radius: var(--radius-lg); padding: 32px 28px; border: 1px solid var(--border-color); transition: border-color .25s ease, box-shadow .25s ease; }
.adv-card:hover { border-color: var(--accent); box-shadow: var(--shadow-gold); }
.adv-icon { width: 56px; height: 56px; display: flex; align-items: center; justify-content: center; background: var(--gradient-gold); border-radius: var(--radius-md); color: #fff; margin-bottom: 16px; }
.adv-card h3 { font-family: var(--font-serif); font-size: 18px; font-weight: 700; margin-bottom: 14px; }
.adv-card li { position: relative; padding-left: 18px; margin-bottom: 7px; font-size: 14px; color: var(--text-secondary); line-height: 1.6; }
.adv-card li::before { content: ''; position: absolute; left: 0; top: 9px; width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }

/* ══════ Founder ══════ */
.founder-profile { display: grid; grid-template-columns: 300px 1fr; gap: 48px; align-items: start; }
.founder-img { border-radius: var(--radius-xl); overflow: hidden; box-shadow: var(--shadow-lg); }
.founder-img img { width: 100%; aspect-ratio: 3/4; object-fit: cover; }
.founder-info h3 { font-family: var(--font-serif); font-size: 1.5rem; font-weight: 700; margin-bottom: 2px; }
.founder-title { font-size: 14px; color: var(--accent); font-weight: 600; margin-bottom: 20px; }
.founder-info > p { font-size: 15px; color: var(--text-secondary); line-height: 1.9; margin-bottom: 12px; }
.founder-quote { border-left: 4px solid var(--accent); padding: 16px 24px; margin: 24px 0; background: rgba(200,150,62,.05); border-radius: 0 var(--radius-md) var(--radius-md) 0; }
.founder-quote p { font-family: var(--font-serif); font-size: 1.1rem; color: var(--primary); font-style: italic; margin: 0; }
.founder-gallery { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-top: 40px; }
.founder-gallery img { width: 100%; aspect-ratio: 4/3; object-fit: cover; border-radius: var(--radius-md); transition: transform .25s ease; }
.founder-gallery img:hover { transform: scale(1.04); }
.founder-divider { border: none; border-top: 1px solid var(--border-color); margin: 64px 0 48px; }

/* ══════ Business ══════ */
.biz-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 28px; }
.biz-card { background: var(--bg-light); border-radius: var(--radius-xl); padding: 36px; border: 1px solid var(--border-color); }
.biz-num { font-family: var(--font-serif); font-size: 3rem; font-weight: 800; color: var(--accent); opacity: .2; line-height: 1; margin-bottom: -12px; }
.biz-card h3 { font-family: var(--font-serif); font-size: 1.25rem; font-weight: 700; margin-bottom: 4px; }
.biz-subtitle { font-size: 13px; color: var(--accent); font-weight: 600; margin-bottom: 18px; }
.biz-card li { position: relative; padding-left: 18px; margin-bottom: 7px; font-size: 14px; color: var(--text-secondary); line-height: 1.6; }
.biz-card li::before { content: ''; position: absolute; left: 0; top: 9px; width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }
.biz-card img { width: 100%; border-radius: var(--radius-md); margin: 16px 0; }
.biz-note { font-size: 12px; color: var(--text-light); text-align: center; margin-top: 4px; }
.biz-highlight { font-size: 14px; color: var(--text-primary); font-weight: 600; line-height: 1.7; margin-top: 12px; }

/* ══════ Honors ══════ */
.honors-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 24px; }
.honors-grid figure { background: var(--bg-white); border-radius: var(--radius-lg); padding: 24px; box-shadow: var(--shadow-sm); transition: box-shadow .25s ease; cursor: pointer; }
.honors-grid figure:hover { box-shadow: var(--shadow-md); }
.honors-grid figure:focus-visible { outline: none; box-shadow: var(--shadow-focus); }
.honors-grid img { width: 100%; aspect-ratio: 3/4; object-fit: contain; }

/* ══════ Quality ══════ */
.quality-card { display: grid; grid-template-columns: 1fr 1fr; gap: 48px; align-items: center; background: var(--gradient-hero); border-radius: var(--radius-xl); padding: 52px; color: #fff; }
.quality-card h3 { font-family: var(--font-serif); font-size: 1.75rem; font-weight: 800; margin-bottom: 16px; }
.quality-card p { font-size: 15px; color: rgba(255,255,255,.7); line-height: 1.85; }
.quality-card img { width: 100%; border-radius: var(--radius-lg); }

/* ══════ Closing ══════ */
.section-closing { position: relative; overflow: hidden; }
.section-closing::before { content: ''; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); width: 600px; height: 600px; background: radial-gradient(circle, rgba(200,150,62,.04), transparent 70%); pointer-events: none; }
.closing { max-width: 720px; margin: 0 auto; text-align: center; position: relative; z-index: 1; }
.closing h2 { font-family: var(--font-serif); font-size: clamp(1.6rem, 3.5vw, 2rem); font-weight: 800; margin-bottom: 28px; }
.closing p { font-size: 16px; color: var(--text-secondary); line-height: 2; margin-bottom: 14px; }
.closing-lead { font-size: 1.15rem; color: var(--text-primary) !important; font-weight: 500; }
.closing-cta { font-size: 1.15rem; font-weight: 700; color: var(--accent-dark) !important; margin-top: 28px !important; }

/* ══════ Focus ══════ */
.card:focus-visible,
.adv-card:focus-visible,
.biz-card:focus-visible,
.honors-grid figure:focus-visible { outline: none; box-shadow: var(--shadow-focus); }

/* ══════ Reveal ══════ */
.reveal { opacity: 0; transform: translateY(28px); transition: opacity .6s ease, transform .6s ease; }
.reveal.visible { opacity: 1; transform: translateY(0); }
.d1 { transition-delay: .1s; }
.d2 { transition-delay: .2s; }
.d3 { transition-delay: .3s; }

/* ══════ Lightbox ══════ */
.cert-lightbox { position: fixed; inset: 0; z-index: 10000; background: rgba(0,0,0,.88); display: flex; align-items: center; justify-content: center; padding: 40px; }
.cert-lightbox img { max-width: 90vw; max-height: 90vh; object-fit: contain; border-radius: var(--radius-lg); }
.cert-lightbox-close { position: absolute; top: 20px; right: 24px; width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; background: rgba(255,255,255,.12); border: none; border-radius: 50%; color: #fff; font-size: 28px; cursor: pointer; transition: background .2s; }
.cert-lightbox-close:hover { background: rgba(255,255,255,.22); }

/* ══════ Responsive ══════ */
@media (max-width: 1023px) {
  .card-grid.cols-3 { grid-template-columns: repeat(2, 1fr); }
  .adv-grid { grid-template-columns: repeat(2, 1fr); }
  .biz-grid { grid-template-columns: 1fr; }
}
@media (max-width: 767px) {
  .hero { padding: 32px 0 56px; }
  .hero h1 { margin: 16px 0 12px; }
  .hero-sub { font-size: var(--text-lg); }
  .hero-desc { font-size: var(--text-sm); }
  .section { padding: var(--section-gap-mobile) 0; }
  .stats-bar { flex-wrap: wrap; justify-content: center; }
  .stat-item { flex: 0 0 auto; min-width: 80px; }
  .intro-grid { grid-template-columns: 1fr; gap: 32px; }
  .intro-img { order: -1; }
  .card-grid.cols-3 { grid-template-columns: 1fr; }
  .img-pair { grid-template-columns: 1fr; }
  .dark-card { padding: 28px 24px; }
  .adv-grid { grid-template-columns: 1fr; }
  .founder-profile { grid-template-columns: 1fr; gap: 24px; }
  .founder-img { max-width: 240px; }
  .founder-gallery { grid-template-columns: repeat(2, 1fr); }
  .honors-grid { grid-template-columns: repeat(2, 1fr); }
  .quality-card { grid-template-columns: 1fr; padding: 32px 24px; gap: 24px; }
}
</style>
