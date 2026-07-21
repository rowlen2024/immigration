<template>
  <div class="about-page">
    <div class="container">
      <ProjectBreadcrumb label="公司简介"/>
    </div>

    <!-- ── Hero ── -->
    <section class="hero">
      <div class="container hero-inner">
        <span class="hero-badge">SINCE 2013</span>
        <h1>关于北极星移民</h1>
        <p class="hero-sub">北极星移民，全球身份规划与移民服务团队</p>
        <p class="hero-desc">北极星移民由深圳福瑞特出入境信息咨询有限公司运营，专注海外身份规划、投资移民与配套服务</p>
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
            <p>北极星移民是深圳福瑞特出入境信息咨询有限公司运营的移民服务品牌，成立于<strong>2013年</strong>，总部位于深圳，专注于海外身份规划与全球资产配置等高端移民服务。</p>
            <p>作为移民领域的权威机构，北极星移民汇聚了一支由资深移民顾问、专业项目风控团队、文案律师团队、税务专家及跨国教育顾问组成的精英团队。他们凭借对各国移民政策的深刻解读与丰富的实操经验，为客户精准匹配最优路径。</p>
            <p>作为<strong>美国EB-5投资移民、香港投资移民、土耳其投资移民、巴拿马购房投资移民、美国EB-1A杰出人才移民、还是欧洲黄金签证、加勒比海等各国护照的先驱者</strong>，北极星向来以数据化分析、个性化方案及全程透明化管理，确保成功率与效率的双重保障。同时，我们更提供涵盖海外置业、税务优化、国际教育等增值服务，真正实现"移民+安家+发展"的一站式解决方案。</p>
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
          <span class="section-label">行业创举</span>
          <h2>行业创举引领，树立专业新标杆</h2>
        </div>

        <div class="card-grid cols-3">
          <div class="card reveal">
            <div class="card-icon" v-html="icons.clock" />
            <h3>13年+行业深耕</h3>
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
            <p>成功举办6届大湾区国际教育名校展，，全网粉丝10万+，首创"身份+教育+资产"三维服务体系。</p>
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
            <p>中欧国际工商学院硕士，深耕<strong>16 年国际移民、8 年国际教育</strong>两大领域，见证行业全发展周期，对全球移民政策、海外项目及行业风险拥有深度研判能力。</p>
            <div class="founder-quote"><p>"坚持每年投入大量时间，亲自带领团队飞赴全球各个项目地实地考察。用脚步丈量项目，只荐真项目、只说真话。"</p></div>
            <p>作为<strong>《聪明的大白》《三只小象部落》《葛瑞说择校》</strong>创始人，《聪明的大白》稳居国际教育行业标杆，同时也是<strong>《粤港澳大湾区国际教育名校展》唯一主办方</strong>。</p>
            <p>常年带队全球实地考察项目，坚守初心、以诚为本。以<strong>"诚信专业·使命必达"</strong>为北极星服务准则，团队精工细作、严格把控每一处细节，真诚服务每一位客户。</p>
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
            <p>张霞女士拥有<strong>18 年移民、7 年国际教育</strong>行业经验，深耕跨界领域，精准洞悉移民家庭身份规划与子女海外择校、教育衔接的核心需求。</p>
            <p>2020 年，她创立知名教育 IP<strong>《K妈谈择校》</strong>，同步运营<strong>《K妈聊香港》</strong>视频号，为众多涉外家庭提供专业升学、身份规划指导。</p>
            <p>她坚守品牌服务宗旨，搭建完善的标准化服务体系，统筹四大部门高效协同。与夏竹青先生携手深耕行业，助力北极星移民凭借专业服务与优良口碑，赢得广大客户信赖与认可。</p>
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
            <p>三只小象联合创始人、总经理，《聪明的大白》公众号主编，深耕国际教育领域<strong>10 年</strong>。</p>
            <p>常年实地探访大湾区数百所国际学校，深度对接胡润先生、多位知名校长等行业大咖，手握海量行业资源。</p>
            <p>凭借扎实的专业能力与实战经验，已助力<strong>五千余个家庭</strong>敲定理想院校，同时全力赋能北极星移民客户，一站式解答各类教育难题，以专业务实的态度与真诚服务，赢得广大家庭的认可与信赖。</p>
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
            <p class="biz-highlight">北极星移民板块优势：深耕美国全品类移民赛道，拥有资深美籍律师团队，擅长疑难案、补件案处理，全程定制方案、极速递交，整体获批率行业领先，一站式搞定美国绿卡及身份落地规划。</p>
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
            <p class="biz-highlight">北极星移民板块优势：业内成功率高，熟知港府审批细则，精准优化文书，覆盖高才/优才/专才/进修/投资全渠道，续签稳、永居成功率高，配套香港升学、开户、落地全套服务。</p>
          </div>

          <!-- 03 巴拿马移民 -->
          <div class="biz-card reveal">
            <div class="biz-num">03</div>
            <h3>巴拿马移民</h3>
            <p class="biz-subtitle">购房移民 · 快速永居</p>
            <ul>
              <li>购房移民：购置 30 万美金房产；审批快，三代移民，移民监宽松，美元资产，美国 E2 签证跳板</li>
            </ul>
            <p class="biz-highlight">北极星移民板块优势：拥有当地一手合规房源资源，全程极速审批、无隐形消费，精通美元资产配置与 E2 签证联动规划，轻松实现身份+资产双布局。</p>
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
            <p class="biz-note">▲北极星移民＆星途ASTER葡萄牙基金合作授权仪式</p>
            <p class="biz-highlight">北极星移民板块优势：严选欧洲合规优质项目与房源/基金渠道，规避政策风险，熟悉各国移民监、入籍政策，定制宜居、养老、欧盟护照专属落地方案。</p>
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
            <p class="biz-highlight">北极星移民板块优势：专注全球一步到位护照项目，全部为政府合规渠道，先批后投、风险可控，精通资产隔离、税务筹划、名校跳板规划，全程线上办理、高效省心。</p>
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

useSeo({
  title: '关于我们',
  description: '北极星移民成立于2013年，专注海外身份规划与全球资产配置。13年+行业经验，2万+成功案例，成功率96.67%。',
  breadcrumbLabel: '关于我们',
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
    title: '专业的文案团队',
    points: ['文案团队平均8年以上行业经验', '深度挖掘客户的背景情况', '20000+案子的申请经历'],
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
    points: ['全部符合各国移民局精准文件', '申请投资类移民成功率100%', '身份办理通过率96.67%+'],
  },
  {
    icon: getIconSvg('award', 28),
    title: '优质的后端服务',
    points: ['跟踪移民项目的后端服务', '完善贴心的异国安家服务', '5年+的持续跟进服务'],
  },
  {
    icon: getIconSvg('landmark', 28),
    title: '跨境财税与家庭配套团队',
    points: ['精通各国跨境税务、资产合规及海外教育规则', '梳理资产规划、规避财税风险', '对接院校资源，协助子女入学', '一站式守护资产安全，解决全家落地难题'],
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
