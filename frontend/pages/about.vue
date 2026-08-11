<template>
  <div class="about-page">
    <div class="container">
      <ProjectBreadcrumb :items="breadcrumbs" />
    </div>

    <!-- ── Hero ── -->
    <section class="hero">
      <div class="container hero-inner">
        <span class="hero-badge">SINCE 2013</span>
        <h1>领航全球<span class="dot">·</span>智启未来</h1>
        <p class="hero-sub">北极星移民，您值得信赖的国际化伙伴</p>
        <p class="hero-desc">2013 年至今持续运营移民业务，累计服务 20000 + 组家庭，提供可量化、可追溯的海外身份、资产、教育一体化落地服务</p>
        <NuxtLink to="/contact" class="btn-primary hero-cta">免费咨询</NuxtLink>
      </div>
    </section>

    <!-- ── Stats ── -->
    <section class="stats-section">
      <div class="container">
        <dl class="stats-bar">
          <div class="stat-item reveal"><dt class="stat-lbl">行业经验</dt><dd class="stat-num">13<span class="stat-unit">年+</span></dd></div>
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
            <p>北极星移民（深圳福瑞特出入境信息咨询有限公司运营品牌），成立于<strong>2013年4月</strong>，总部位于深圳，专注于海外身份规划与全球资产配置等高端移民服务。</p>
            <p>公司持有广东省因私出入境移民协会会员单位、5A级诚信企业、美国IIUSA正式会员、美国松州区域中心官方授权合作伙伴、ASTER葡萄牙基金全球战略合作伙伴、巴拿马官方移民服务中心全球战略合作伙伴、格林纳达国家度假村项目铂金级合作伙伴、JWP地产全球指定合作伙伴、深圳市跨境电子商务协会副会长单位、中欧深圳校友会副会长单位等12项行业资质与官方授权认证。</p>
            <p>业务范围覆盖<strong>美国EB-5投资移民</strong>、<strong>香港投资移民</strong>、<strong>巴拿马购房投资移民</strong>、<strong>美国EB-1A杰出人才移民</strong>、<strong>欧洲黄金签证</strong>、<strong>加勒比海各国护照</strong>等品类，以数据化分析、个性化方案及全程透明化管理确保办理效率。同时提供海外置业、税务优化、国际教育等增值服务，实现「移民+安家+发展」的一站式解决方案。</p>
          </div>
          <div class="intro-img reveal">
            <img :src="`/static/about/images/intro.jpg`" alt="北极星移民团队" width="800" height="600" loading="lazy" />
          </div>
        </div>
      </div>
    </section>

    <!-- ── Milestones ── -->
    <section class="section section-alt">
      <div class="container">
        <div class="section-header reveal">
          <span class="section-label">企业核心实力</span>
          <h2>行业先发优势 | 交付实绩</h2>
        </div>

        <div class="card-grid cols-3">
          <div class="card reveal">
            <div class="card-icon" v-html="icons.clock" />
            <h3>13 年持续经营沉淀</h3>
            <p>核心管理层每年实地走访全球移民项目，累计实地调研 60 + 国家及地区，累计服务2万+ 组签约客户，留存完整可查申请档案。</p>
          </div>
          <div class="card reveal d1">
            <div class="card-icon" v-html="icons.diamond" />
            <h3>巴拿马投资永居首批服务机构</h3>
            <p>国内最早开展巴拿马合格投资者签证推广的机构之一，累计交付 15 组完整获批家庭案例，董事长亲赴巴拿马与移民局、法案修订官员面对面沟通政策细则。</p>
          </div>
          <div class="card reveal d2">
            <div class="card-icon" v-html="icons.shield" />
            <h3>EB5 全流程落地数据</h3>
            <p>国内具备 EB5 完整周期办理资质企业，累计 290 + 组 EB-5 获批案卷；保持中国大陆申请人 EB-1A 8 天加急获批行业记录；完成 70 周岁以上高龄申请人 EB-1A 完整获批案例。</p>
          </div>
          <div class="card reveal">
            <div class="card-icon" v-html="icons.check" />
            <h3>香港全通道稳定交付</h3>
            <p>同步运营资本投资者入境计划、高才通 A/B/C 类、优才、专才、进修 5 类香港身份通道；香港各类签证续签递交案例 99% 通过入境处审核；香港投资移民累计交付 200 + 组完整转永居案例。</p>
          </div>
          <div class="card reveal d1">
            <div class="card-icon" v-html="icons.trending" />
            <h3>全业务线通过率数据</h3>
            <p>全品类移民综合申请通过率 96.67%；所有投资类移民项目递交申请获批率 99%；非投资类人才签证整体通过率 95% 以上。</p>
          </div>
          <div class="card reveal d2">
            <div class="card-icon" v-html="icons.book" />
            <h3>旗下教育品牌"三只小象"</h3>
            <p>独立主办 6 届大湾区国际教育名校展；全网矩阵账号累计粉丝总量 10 万 +；搭建行业独有「身份规划 + 国际择校 + 跨境资产」三维配套服务体系。</p>
          </div>
        </div>

        <div class="img-pair reveal">
          <figure>
            <img :src="`/static/about/images/1778568651667806803.jpg`" alt="北极星移民成功案例墙" width="800" height="500" loading="lazy" />
            <figcaption>▲北极星移民成功案例墙</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/1778568704715310459.jpg`" alt="巴拿马论坛合影" width="800" height="500" loading="lazy" />
            <figcaption>▲与巴拿马合格投资者法案修订人RICHARD THOMAS在投资者永居计划巡回论坛合影</figcaption>
          </figure>
        </div>

        <div class="dark-card reveal">
          <h4>战略孵化教育子品牌"三只小象"</h4>
          <ul>
            <li>独家冠名6届大湾区规模最大的国际教育展</li>
            <li>拥有大湾区最具影响力的全媒体平台，全网粉丝10万＋</li>
            <li>全网唯一一家探校广州民心学校平台，合作学校涵盖哈罗国际、哈罗港人、华附国际部、广州暨大、南山国际、蛇口国际、科爱赛、南山国王、新哲、预科、本一幼儿园、金生幼儿园等大湾区众多知名国际学校</li>
            <li>首创"身份+教育+资产"三维服务体系，直击家庭核心诉求，让移民真正成为高净值家庭进阶的跳板</li>
            <li><strong>旗下知名账号：</strong>《聪明的大白》《葛瑞说择校》《三只小象部落》等公众号，小红书《vivi帮你择校》《K妈谈择校》等十几个账号</li>
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
        <el-carousel
          class="honors-carousel reveal"
          :type="honorCarouselType"
          height="420px"
          trigger="click"
          arrow="always"
          indicator-position="outside"
          :autoplay="true"
          :interval="4000"
          :loop="true"
          :pause-on-hover="true"
        >
          <el-carousel-item v-for="(h, i) in honorImgs" :key="h.src">
            <figure tabindex="0" role="button" :aria-label="`查看${h.alt}大图`" @click="selectedCert = i" @keydown.enter="selectedCert = i" @keydown.space.prevent="selectedCert = i">
              <img :src="`/static/about/images/${h.src}`" :alt="h.alt" width="600" height="800" loading="lazy" />
            </figure>
          </el-carousel-item>
        </el-carousel>
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
            <img :src="`/static/about/images/夏竹青.jpg`" alt="创始人夏竹青先生" width="600" height="800" loading="lazy" />
          </div>
          <div class="founder-info">
            <h3>夏竹青 先生</h3>
            <p class="founder-title">创始人兼董事长</p>
            <p>中欧国际工商学院硕士，胡润百富“国际化教育影响力人物”，累计<strong>16 年移民行业全职从业经验、8 年国际教育服务运营经验</strong>，完整覆盖国内移民行业 2013-2026 全发展周期，具备全球移民政策、海外投资项目、行业风险完整研判履历。</p>
            <div class="founder-quote"><p>"坚持每年带队出境实地走访全球移民项目场地，实地核验项目资质、施工进度、回款机制；对外服务标准：仅推送完成实地尽调项目，对外公示项目全部客观政策与收益数据。"</p></div>
            <p>作为<strong>《聪明的大白》《三只小象部落》《葛瑞说择校》</strong>创始人，《聪明的大白》稳居国际教育行业标杆，同时也是<strong>《粤港澳大湾区国际教育名校展》唯一主办方</strong>。</p>
            <p>常年带队全球实地考察项目，确立企业「<strong>项目实地尽调</strong> + <strong>海外本土律师双审核</strong>」项目风控标准，全流程标准化管控申请节点，企业运营至今 0 客户诉讼、0 服务纠纷。</p>
          </div>
        </div>

        <div class="founder-gallery reveal">
          <figure>
            <img :src="`/static/about/images/和美国移民局EB5办公室前副主任合影.jpg`" alt="和美国移民局EB5办公室前副主任合影" width="400" height="300" loading="lazy" />
            <figcaption>和美国移民局EB5办公室前副主任合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和EB5鼻祖律师张律师合影.jpg`" alt="和EB5鼻祖律师张律师合影" width="400" height="300" loading="lazy" />
            <figcaption>和EB5鼻祖律师张律师合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和巴拿马移民局官员合影.jpg`" alt="和巴拿马移民局官员合影" width="400" height="300" loading="lazy" />
            <figcaption>和巴拿马移民局官员合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和巴拿马移民法案修订人合影.jpg`" alt="和巴拿马移民法案修订人合影" width="400" height="300" loading="lazy" />
            <figcaption>和巴拿马移民法案修订人合影</figcaption>
          </figure>
        </div>

        <hr class="founder-divider" />

        <div class="founder-profile reveal">
          <div class="founder-img">
            <img :src="`/static/about/images/张霞.jpg`" alt="联合创始人张霞女士" width="600" height="800" loading="lazy" />
          </div>
          <div class="founder-info">
            <h3>张霞 女士</h3>
            <p class="founder-title">联合创始人 · 运营总监</p>
            <p>张霞女士拥有<strong>18 年移民行业实操经验、7 年国际教育行业经验</strong>，深耕移民身份规划 + 子女海外择校双赛道，精准匹配家庭身份、教育两类硬性需求。</p>
            <p>2020 年创立垂直教育 IP<strong>《K妈谈择校》</strong>，同步运营视频号<strong>《K妈聊香港》</strong>，累计服务上千组涉外家庭升学、国际身份规划案卷。</p>
            <p>搭建公司四大职能部门标准化协作流程，统筹顾问、文案、风控、售后全链路运营；与董事长夏竹青先生协同搭建全流程服务体系，经营周期客户 0 诉讼、0 纠纷记录。</p>
          </div>
        </div>

        <div class="founder-gallery reveal">
          <figure>
            <img :src="`/static/about/images/和美国佛蒙特州前州长合影.jpg`" alt="和美国佛蒙特州前州长合影" width="400" height="300" loading="lazy" />
            <figcaption>和美国佛蒙特州前州长合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和美国前议员合影.jpg`" alt="和美国前议员合影" width="400" height="300" loading="lazy" />
            <figcaption>和美国前议员合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和巴拿马移民法案修订人合影-张霞.jpg`" alt="和巴拿马移民法案修订人合影" width="400" height="300" loading="lazy" />
            <figcaption>和巴拿马移民法案修订人合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和葡萄牙最大的侓师事务所合伙人合影-张霞.jpg`" alt="和葡萄牙最大的律师事务所合伙人合影" width="400" height="300" loading="lazy" />
            <figcaption>和葡萄牙最大的律师事务所合伙人合影</figcaption>
          </figure>
        </div>

        <hr class="founder-divider" />

        <!-- 大白老师 -->
        <div class="founder-profile reveal">
          <div class="founder-img">
            <img :src="`/static/about/images/大白老师.jpg`" alt="大白老师" width="600" height="800" loading="lazy" />
          </div>
          <div class="founder-info">
            <h3>大白老师</h3>
            <p class="founder-title">三只小象联合创始人 · 总经理</p>
            <p>三只小象联合创始人、总经理，<strong>《聪明的大白》</strong>公众号主编，国际教育领域从业<strong>10年。</strong></p>
            <p>常年实地探访大湾区国际学校，数量逾百所。与胡润先生、多位国际学校校长等行业人士有过交流对接。</p>
            <p>已服务五千余个家庭完成国际学校择校。在北极星移民项目中负责为客户提供教育相关问题咨询与解答。</p>
          </div>
        </div>

        <div class="founder-gallery reveal">
          <figure>
            <img :src="`/static/about/images/和胡润先生合影.jpg`" alt="和胡润先生合影" width="400" height="300" loading="lazy" />
            <figcaption>和胡润先生合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和深外国际部创校校长邬晓莉校长合影.jpg`" alt="和深外国际部创校校长邬晓莉校长合影" width="400" height="300" loading="lazy" />
            <figcaption>和深外国际部创校校长邬晓莉校长合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和维多利亚前总校长孔美琪博士合影.jpg`" alt="和维多利亚前总校长孔美琪博士合影" width="400" height="300" loading="lazy" />
            <figcaption>和维多利亚前总校长孔美琪博士合影</figcaption>
          </figure>
          <figure>
            <img :src="`/static/about/images/和哈罗港人子弟学校创校校长合影.jpg`" alt="和哈罗港人子弟学校创校校长合影" width="400" height="300" loading="lazy" />
            <figcaption>和哈罗港人子弟学校创校校长合影</figcaption>
          </figure>
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
          <!-- 01 美国移民 -->
          <div class="biz-card reveal">
            <div class="biz-num">01</div>
            <h3>美国移民</h3>
            <p class="biz-subtitle">人才投资通道全覆盖</p>
            <ul>
              <li><strong>EB-5 投资移民：</strong>投资 80 万美金，创造 10 个岗位；全家绿卡，无语言学历要求，5 年可入籍，教育福利好</li>
              <li><strong>EB-1A 杰出人才：</strong>满足 3 项行业标准；零投资、无雇主，审批快，适合高端人才</li>
              <li><strong>NIW 国家利益豁免：</strong>高学历 / 专业技能、行业贡献达标；门槛低于 EB-1A，无需雇主担保</li>
              <li><strong>EB-1C 跨国高管：</strong>中美关联公司，海外高管满 1 年；无劳工证、无配额、审批快，适合企业主 / 高管，可先 L-1A 再转绿卡</li>
            </ul>
            <p class="biz-highlight">北极星移民板块优势：覆盖美国全品类移民通道，配备美籍律师团队，处理疑难案件及补件案件；提供全程方案定制及材料递交服务。服务范围涵盖美国绿卡申请及后续身份落地规划。。</p>
          </div>

          <!-- 02 中国香港身份 -->
          <div class="biz-card reveal d1">
            <div class="biz-num">02</div>
            <h3>中国香港身份</h3>
            <p class="biz-subtitle">多通道 · 高成功率</p>
            <ul>
              <li><strong>资本投资者入境计划：</strong>投资 3000 万港币；无强制居住要求，7 年转特区护照，低税率，升学优势突出</li>
              <li><strong>高才通：</strong>高薪 / 百强名校背景；零投资、下签快，可自由工作定居</li>
              <li><strong>优才计划：</strong>满足 12 项目中 6 项标准；零投资，续签灵活，适配各类专业人才</li>
              <li><strong>专才计划：</strong>雇主担保 / 本地创业；无行业配额限制，审批高效，配偶可工作、子女享公立教育</li>
              <li><strong>进修移民：</strong>就读全日制课程；入门门槛低，毕业获留港签证，7 年可转永居</li>
            </ul>
            <p class="biz-highlight">北极星移民板块优势：覆盖高才、优才、专才、进修、投资等全部申请通道，掌握香港入境处审批细则及文书优化方案。配套服务涵盖续签指导、永居申请、香港升学、银行开户及落地支持。</p>
          </div>

          <!-- 03 巴拿马移民 -->
          <div class="biz-card reveal">
            <div class="biz-num">03</div>
            <h3>巴拿马移民</h3>
            <p class="biz-subtitle">购房移民 · 快速永居</p>
            <ul>
              <li>购房移民：购置 30 万美金房产；审批快，三代移民，移民监宽松，美元资产，美国 E2 签证跳板</li>
            </ul>
            <p class="biz-highlight">北极星移民板块优势：提供当地一手合规房源资源，全程办理无附加收费。服务范围涵盖美元资产配置与E2签证联动规划，协助完成身份与资产同步布局。</p>
          </div>

          <!-- 04 欧洲（黄金签证 / 欧盟身份） -->
          <div class="biz-card reveal d1">
            <div class="biz-num">04</div>
            <h3>欧洲（黄金签证 / 欧盟身份）</h3>
            <p class="biz-subtitle">希腊 · 葡萄牙 · 马耳他</p>
            <ul>
              <li><strong>希腊：</strong>购房 25 万欧元；无移民监，快速拿永居，畅行申根区，宜居养老优选</li>
              <li><strong>葡萄牙：</strong>50 万欧元基金；移民监宽松，满 5 年可入籍拿欧盟护照</li>
              <li><strong>马耳他：</strong>捐款 + 租房组合；四代移民，英语环境，集欧盟、申根、欧元、英联邦身份于一体</li>
            </ul>
            <!--
            <p class="biz-note">▲北极星移民＆星途ASTER葡萄牙基金合作授权仪式</p>
            -->
            <p class="biz-highlight">北极星移民板块优势：筛选欧洲合规项目及房源/基金渠道，规避政策风险；熟悉各国移民监及入籍政策要求，提供宜居、养老、欧盟护照等专属落地方案。</p>
          </div>

          <!-- 05 投资入籍 -->
          <div class="biz-card reveal">
            <div class="biz-num">05</div>
            <h3>投资入籍</h3>
            <p class="biz-subtitle">一步到位护照项目</p>
            <ul>
              <li><strong>土耳其：</strong>购房 40 万美金 / 基金 50 万美金（持有 3 年）；无移民监，大国身份，美国 E2 优质跳板，子女可走国内名校国际生通道</li>
              <li><strong>圣基茨和尼维斯：</strong>老牌合规项目，零登陆零移民监，免税体系，英联邦身份，免签国家众多</li>
              <li><strong>圣多美和普林西比：</strong>先批后投，入门成本低，未参与 CRS，可衔接葡语国家身份</li>
              <li><strong>圣卢西亚：</strong>政策灵活，无居住要求，税务友好，出行便利</li>
              <li><strong>格林纳达：</strong>加勒比独有美国 E2 资质，免签中国大陆，支持四代移民，税务优势明显</li>
              <li><strong>安提瓜和巴布达：</strong>代同行性价比高，免签实力强，零移民监、无个税</li>
              <li><strong>瓦努阿图：</strong>全球审批速度最快，全程无需登陆，非 CRS 国家，纯免税</li>
              <li><strong>多米尼克：</strong>加勒比低价优选，支持四代移民，无居住要求，资产隐私性强</li>
            </ul>
            <p class="biz-highlight">北极星移民板块优势：专注于全球一步到位护照项目，全部采用政府合规渠道，以先获批后投资模式为主。业务范围涵盖资产隔离、税务筹划、国际生通道规划等，全程可线上办理。</p>
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
import { buildAboutPageJsonLd, toJsonLdScripts } from '~/utils/jsonld'
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'

/** 关于我们页面的稳定面包屑路径。 */
const breadcrumbs = [
  HOME_BREADCRUMB,
  { label: '关于我们', href: '/about' },
]

useSeo({
  title: '关于我们',
  description: '北极星移民成立于2013年，专注海外身份规划与全球资产配置。13年+行业经验，2万+成功案例，成功率96.67%。',
  breadcrumbs,
})

useHead(() => ({
  script: toJsonLdScripts(buildAboutPageJsonLd(
    '关于北极星移民',
    '北极星移民（深圳福瑞特出入境信息咨询有限公司）成立于2013年，总部位于深圳，专注海外身份规划与全球资产配置等高端移民服务。覆盖美国EB-5/EB-1A、香港投资移民、巴拿马购房移民、欧洲黄金签证及加勒比护照项目，成功率96.67%。',
    {
      '@type': 'Organization',
      name: '北极星移民',
      foundingDate: '2013',
      description: '专业投资移民服务机构，13年+行业深耕，2万+成功案例',
    },
  )),
}))

onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => entries.forEach((e) => { if (e.isIntersecting) e.target.classList.add('visible') }),
    { threshold: 0.15 },
  )
  document.querySelectorAll('.reveal').forEach((el) => observer.observe(el))

  honorCarouselMedia = window.matchMedia('(max-width: 767px)')
  updateHonorCarouselType(honorCarouselMedia)
  honorCarouselMedia.addEventListener('change', updateHonorCarouselType)
})

const honorImgs = [
  { src: '1778637606374142161.jpg', alt: '北极星移民企业荣誉证书 — 行业资质认证' },
  { src: '1778637708046637702.jpg', alt: '北极星移民企业荣誉证书 — 优秀服务机构' },
  { src: '1778637719093755072.jpg', alt: '北极星移民企业荣誉证书 — 诚信经营示范单位' },
  { src: '1778637728035122831.jpg', alt: '北极星移民企业荣誉证书 — 行业标杆品牌' },
  { src: '20260721140058.jpg', alt: '北极星移民企业荣誉证书 — 副会长单位' },
]

const selectedCert = ref<number | null>(null)
const honorCarouselType = ref<'card' | ''>('card')
let honorCarouselMedia: MediaQueryList | null = null

function updateHonorCarouselType(event: MediaQueryList | MediaQueryListEvent) {
  honorCarouselType.value = event.matches ? '' : 'card'
}

onUnmounted(() => {
  honorCarouselMedia?.removeEventListener('change', updateHonorCarouselType)
})

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
    title: '标准化文案团队',
    points: ['文案专员人均在岗从业时长 8 年以上', '案卷累计总交付量 1000 + 套完整申请材料', '固定执行申请人背景全维度材料挖掘文书标准流程', '单个文案平均经手个案超过2000+'],
  },
  {
    icon: getIconSvg('users', 28),
    title: '持证顾问团队',
    points: ['顾问团队平均从业时间10年以上', '固定执行 1 对 1 专属顾问终身跟进服务机制', '全员熟记全球各国移民法案细则、材料红线、拒签风险节点'],
  },
  {
    icon: getIconSvg('target', 28),
    title: '风险把控的项目团队',
    points: ['项目风控团队平均行业从业年限 16 年', '实地尽调全球 60 + 海外移民投资项目', '合作项目自合作至今无投资风险纠纷、无客户诉讼记录'],
  },
  {
    icon: getIconSvg('trending-up', 28),
    title: '可追溯申请通过率',
    points: ['全部案卷严格匹配各国移民局材料清单与审核标准', '投资类移民案卷获批率 99%', '综合身份办理整体通过率 96.67% 以上'],
  },
  {
    icon: getIconSvg('award', 28),
    title: '长期后端落地配套服务',
    points: ['设立移民落地专属售后跟踪岗位', '配套海外居住、开户、入学标准化安家服务包', '申请人获批后提供最低 5 年持续售后跟进服务'],
  },
  {
    icon: getIconSvg('landmark', 28),
    title: '跨境财税 & 教育配套团队',
    points: ['全员熟悉多国跨境税务法案、海外资产合规申报、国际学校入学政策', '提供家庭跨境资产梳理、税务风险规避标准化方案', '直连大湾区国际院校招生渠道，全程代办子女入学申请'],
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
.section-label { display: inline-block; font-size: 18px; font-weight: 700; text-transform: uppercase; letter-spacing: .14em; color: var(--accent); margin-bottom: 8px; }
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
.founder-gallery figure { margin: 0; }
.founder-gallery img { width: 100%; aspect-ratio: 4/3; object-fit: cover; border-radius: var(--radius-md); transition: transform .25s ease; }
.founder-gallery img:hover { transform: scale(1.04); }
.founder-gallery figcaption { text-align: center; font-size: 12px; color: var(--text-light); margin-top: 8px; line-height: 1.5; }
.founder-divider { border: none; border-top: 1px solid var(--border-color); margin: 64px 0 48px; }

/* ══════ Business ══════ */
.biz-grid { column-count: 2; column-gap: 28px; }
.biz-card { break-inside: avoid; margin-bottom: 28px; background: var(--bg-light); border-radius: var(--radius-xl); padding: 36px; border: 1px solid var(--border-color); }
.biz-num { font-family: var(--font-serif); font-size: 3rem; font-weight: 800; color: var(--accent); opacity: .2; line-height: 1; margin-bottom: -12px; }
.biz-card h3 { font-family: var(--font-serif); font-size: 1.25rem; font-weight: 700; margin-bottom: 4px; }
.biz-subtitle { font-size: 13px; color: var(--accent); font-weight: 600; margin-bottom: 18px; }
.biz-card li { position: relative; padding-left: 18px; margin-bottom: 7px; font-size: 14px; color: var(--text-secondary); line-height: 1.6; }
.biz-card li::before { content: ''; position: absolute; left: 0; top: 9px; width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }
.biz-card img { width: 100%; border-radius: var(--radius-md); margin: 16px 0; }
.biz-note { font-size: 12px; color: var(--text-light); text-align: center; margin-top: 4px; }
.biz-highlight { font-size: 14px; color: var(--text-primary); font-weight: 600; line-height: 1.7; margin-top: 12px; }

/* ══════ Honors ══════ */
.honors-carousel { margin: 0 auto; }
.honors-carousel figure {
  height: 100%;
  margin: 0 12px;
  padding: 24px;
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: box-shadow .25s ease, transform .25s ease;
}
.honors-carousel figure:hover { box-shadow: var(--shadow-md); transform: translateY(-2px); }
.honors-carousel figure:focus-visible { outline: none; box-shadow: var(--shadow-focus); }
.honors-carousel img { width: 100%; height: 100%; object-fit: contain; }
.honors-carousel :deep(.el-carousel__item) { padding: 10px 0; }
.honors-carousel :deep(.el-carousel__item--card) { opacity: .55; }
.honors-carousel :deep(.el-carousel__item--card.is-active) { opacity: 1; }
.honors-carousel :deep(.el-carousel__arrow) {
  width: 44px;
  height: 44px;
  color: #fff;
  background: rgba(200, 150, 62, .52);
  box-shadow: 0 4px 14px rgba(200, 150, 62, .18);
}
.honors-carousel :deep(.el-carousel__arrow:hover) { background: rgba(200, 150, 62, .78); }
.honors-carousel :deep(.el-carousel__button) { background-color: var(--accent); }

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
.honors-carousel figure:focus-visible { outline: none; box-shadow: var(--shadow-focus); }

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
  .biz-grid { column-count: 1; }
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
  .honors-carousel { margin: 0 -12px; }
  .honors-carousel :deep(.el-carousel__container) { height: 360px !important; }
  .honors-carousel figure { margin: 0 6px; padding: 16px; }
  .honors-carousel :deep(.el-carousel__arrow) { width: 38px; height: 38px; }
  .quality-card { grid-template-columns: 1fr; padding: 32px 24px; gap: 24px; }
}
</style>
