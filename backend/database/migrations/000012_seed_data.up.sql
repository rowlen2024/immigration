SET NAMES utf8mb4;
-- Seed admin user (password: admin123, bcrypt cost=12)
INSERT INTO `users` (`username`, `password_hash`, `display_name`, `role`, `status`)
VALUES ('admin', '$2a$12$jUfCLKTrO8VrOiQS8nrDr.ItDI5UStTzKgZq8VKPjM2zozjnUanZC', '超级管理员', 'admin', 1),
       ('editor', '$2a$12$jUfCLKTrO8VrOiQS8nrDr.ItDI5UStTzKgZq8VKPjM2zozjnUanZC', '内容编辑', 'editor', 1),
       ('viewer', '$2a$12$jUfCLKTrO8VrOiQS8nrDr.ItDI5UStTzKgZq8VKPjM2zozjnUanZC', '只读用户', 'viewer', 1);

-- Seed projects
INSERT INTO `projects` (`slug`, `name`, `country`, `flag_emoji`, `tagline`, `investment_amount`, `investment_value`, `investment_currency`, `processing_period`, `target_crowd`, `overview_title`, `overview_text`, `policy_title`, `policy_text`, `costs_total`, `costs_note`, `cta_text`, `hero_title`, `hero_desc`, `hero_gradient`, `sort_order`, `status`)
VALUES
('eb5', '美国EB-5投资移民', '美国', '🇺🇸', '80万美元投资，全家获美国绿卡', '80万美元', 800000.00, 'USD', '约2-3年', '高净值人士、企业家、希望子女接受美国教育的家庭', 'EB-5投资移民概述', 'EB-5投资移民项目（Employment-Based Fifth Preference）由美国国会于1990年设立，旨在通过外国投资刺激美国经济并创造就业机会。2022年《EB-5改革与诚信法案》（RIA）对项目进行了重大改革，将最低投资额提高至80万美元（TEA地区），并加强了项目监管和诚信措施。', 'EB-5最新政策解读', '2022年3月15日生效的《EB-5改革与诚信法案》对EB-5项目进行了全面改革。新法案保留了TEA地区80万美元和非TEA地区105万美元的投资门槛，为乡村、高失业区和基础设施项目预留了32%的签证配额。', '约85-90万美元', '含投资款、管理费、律师费、I-526E申请费等', '立即咨询EB-5项目', '美国EB-5投资移民', '80万美元投资 · 全家绿卡 · 无语言学历要求', 'linear-gradient(135deg, #1a3a5c 0%, #2d5a8e 50%, #c8963e 100%)', 1, 1),
('cies', '香港资本投资者入境计划', '中国香港', '🇭🇰', '3000万港元投资，获香港居留权', '3000万港元', 3870000.00, 'USD', '约6-12个月', '高净值人士、企业家、寻求国际化身份的商业人士', '香港资本投资者入境计划概述', '香港资本投资者入境计划（Capital Investment Entrant Scheme，简称CIES）是香港特别行政区政府推出的投资移民计划。2024年3月1日重新开放申请，新计划的投资门槛提高至3000万港元，且对投资组合有更严格的要求。', 'CIES最新政策解读', '2024年3月1日起，香港CIES新计划正式生效。投资门槛提高至3000万港元，其中至少300万港元需投资于政府指定的"资本投资者入境计划投资组合"。申请人须证明在提出申请前的两年内，一直拥有不少于3000万港元的净资产。', '约3000-3200万港元', '含投资额、申请费、专业服务费等', '立即咨询香港CIES', '香港资本投资者入境计划', '3000万港元投资 · 无语言要求 · 7年后可申请永居', 'linear-gradient(135deg, #c8963e 0%, #d4a94e 50%, #1a3a5c 100%)', 2, 1),
('panama', '巴拿马购房移民', '巴拿马', '🇵🇦', '30万美元购房，获巴拿马永久居留', '30万美元', 300000.00, 'USD', '约3-6个月', '退休人士、投资者、希望获得第二身份的中产家庭', '巴拿马购房移民概述', '巴拿马购房移民项目（Panama Friendly Nations Visa - Real Estate Investment）允许外国人通过在巴拿马购买价值30万美元以上的房产获得永久居留权。巴拿马使用美元作为流通货币，拥有现代化的基础设施和优惠的税收政策。', '巴拿马购房移民政策解读', '巴拿马政府为吸引外国投资，提供多种移民途径。购房移民是最受欢迎的方式之一，最低投资额为30万美元。获得永久居留权5年后可申请入籍。巴拿马不要求申请人在当地居住，且承认双重国籍。', '约32-35万美元', '含购房款、税费、律师费等', '立即咨询巴拿马移民', '巴拿马购房移民', '30万美元购房 · 3-6个月获批 · 无居住要求', 'linear-gradient(135deg, #1a3a5c 0%, #1a5c3a 50%, #c8963e 100%)', 3, 1);

-- Seed requirements for each project
-- EB5 requirements
INSERT INTO `requirements` (`project_id`, `label`, `is_required`, `sort_order`) VALUES
(1, '年满18周岁', 1, 1),
(1, '投资金额80万美元（TEA地区）或105万美元（非TEA地区）', 1, 2),
(1, '投资资金来源合法', 1, 3),
(1, '创造至少10个全职就业岗位', 1, 4),
(1, '无犯罪记录', 1, 5),
(1, '无语言要求', 0, 6),
(1, '无学历要求', 0, 7);

-- CIES requirements
INSERT INTO `requirements` (`project_id`, `label`, `is_required`, `sort_order`) VALUES
(2, '年满18周岁', 1, 1),
(2, '在提出申请前两年内，拥有不少于3000万港元的净资产', 1, 2),
(2, '投资3000万港元于获许投资资产', 1, 3),
(2, '其中至少300万港元投资于"资本投资者入境计划投资组合"', 1, 4),
(2, '无犯罪记录', 1, 5),
(2, '无语言要求', 0, 6),
(2, '无学历要求', 0, 7);

-- Panama requirements
INSERT INTO `requirements` (`project_id`, `label`, `is_required`, `sort_order`) VALUES
(3, '年满18周岁', 1, 1),
(3, '在巴拿马购买价值30万美元以上的房产', 1, 2),
(3, '无犯罪记录', 1, 3),
(3, '身体健康', 1, 4),
(3, '无语言要求', 0, 5),
(3, '无学历要求', 0, 6);

-- Seed cost items for each project
-- EB5 costs
INSERT INTO `cost_items` (`project_id`, `name`, `amount`, `amount_value`, `amount_currency`, `note`, `sort_order`) VALUES
(1, '投资金额', '80万美元', 800000.00, 'USD', 'TEA地区最低投资额，5年后可返还', 1),
(1, '项目管理费', '5-8万美元', 65000.00, 'USD', '区域中心管理费', 2),
(1, '律师费', '2-3万美元', 25000.00, 'USD', '移民律师服务费', 3),
(1, 'I-526E申请费', '11,160美元', 11160.00, 'USD', 'USCIS申请费', 4),
(1, '体检费', '约500美元', 500.00, 'USD', '指定医院体检', 5),
(1, '翻译公证费', '约2,000美元', 2000.00, 'USD', '文件翻译及公证', 6);

-- CIES costs
INSERT INTO `cost_items` (`project_id`, `name`, `amount`, `amount_value`, `amount_currency`, `note`, `sort_order`) VALUES
(2, '投资金额', '3000万港元', 3870000.00, 'USD', '获许投资资产，7年后可退出', 1),
(2, '投资组合', '300万港元', 38700.00, 'USD', '政府指定投资组合（从3000万中分配）', 2),
(2, '申请费', '约1万港元', 1290.00, 'USD', '入境处申请费用', 3),
(2, '专业服务费', '约15-20万港元', 22600.00, 'USD', '律师、会计师等专业费用', 4),
(2, '体检费', '约2,000港元', 258.00, 'USD', '指定医院体检', 5),
(2, '翻译公证费', '约5,000港元', 645.00, 'USD', '文件翻译及公证', 6);

-- Panama costs
INSERT INTO `cost_items` (`project_id`, `name`, `amount`, `amount_value`, `amount_currency`, `note`, `sort_order`) VALUES
(3, '购房款', '30万美元起', 300000.00, 'USD', '符合移民局要求的房产', 1),
(3, '房产转让税', '约2%', 6000.00, 'USD', '房产价值的2%', 2),
(3, '律师费', '约3,000-5,000美元', 4000.00, 'USD', '房产过户及移民律师', 3),
(3, '移民申请费', '约2,000美元', 2000.00, 'USD', '移民局申请费用', 4),
(3, '翻译公证费', '约1,000美元', 1000.00, 'USD', '文件翻译及公证', 5);

-- Seed timeline phases
-- EB5: 5 phases
INSERT INTO `timeline_phases` (`project_id`, `phase_number`, `phase_name`, `duration`, `title`, `description`, `sort_order`) VALUES
(1, 1, '第一阶段', '1-2个月', '选择项目与签约', '选择EB-5投资项目（区域中心或直投项目），签署投资协议，聘用移民律师，准备资金来源证明文件。', 1),
(1, 2, '第二阶段', '2-3个月', '递交I-526E申请', '向美国移民局（USCIS）递交I-526E移民申请，包括投资证明、资金来源证明、个人背景资料等。获得优先日期（Priority Date）。', 2),
(1, 3, '第三阶段', '12-18个月', 'I-526E审批等待', 'USCIS审核I-526E申请。审理期间可申请工卡和旅行证（I-765/I-131）。审批通过后进入国家签证中心（NVC）阶段。', 3),
(1, 4, '第四阶段', '4-6个月', '领事馆面签或调整身份', '在中国：广州领事馆面签，获得移民签证后入境美国。在美国：递交I-485调整身份申请。入境后获得2年临时绿卡。', 4),
(1, 5, '第五阶段', '6-12个月', '移除条件获永久绿卡', '临时绿卡到期前90天递交I-829申请移除条件。证明投资持续且已创造10个就业岗位。获批后获得10年永久绿卡。', 5);

-- CIES: 5 phases
INSERT INTO `timeline_phases` (`project_id`, `phase_number`, `phase_name`, `duration`, `title`, `description`, `sort_order`) VALUES
(2, 1, '第一阶段', '1-2个月', '资格评估与文件准备', '评估净资产是否符合3000万港元要求，准备资产证明文件（银行证明、审计报告等），聘请香港执业律师。', 1),
(2, 2, '第二阶段', '1-2个月', '递交申请', '向香港入境事务处递交资本投资者入境计划申请，包括净资产证明、投资计划书、个人背景资料。', 2),
(2, 3, '第三阶段', '3-6个月', '原则上批准', '入境处审核申请，发出"原则上批准"通知书。申请人需在6个月内完成投资，并提交投资证明。', 3),
(2, 4, '第四阶段', '1-2个月', '正式批准与入境', '入境处核实投资后发出正式批准，申请人获得24个月逗留期限。此后每次续期24个月。', 4),
(2, 5, '第五阶段', '7年', '连续居住与永久居留', '连续在港居住满7年后，可申请成为香港永久性居民。7年期间需持续满足投资要求。', 5);

-- Panama: 4 phases + E2
INSERT INTO `timeline_phases` (`project_id`, `phase_number`, `phase_name`, `duration`, `title`, `description`, `sort_order`) VALUES
(3, 1, '第一阶段', '1个月', '选房与签约', '赴巴拿马考察房产（或远程选房），签署购房合同，支付定金。委托当地律师办理后续手续。', 1),
(3, 2, '第二阶段', '1-2个月', '房产过户', '完成房产交易，支付房款余额，办理房产证（Title of Property）。缴纳房产转让税和注册费。', 2),
(3, 3, '第三阶段', '2-3个月', '递交移民申请', '向巴拿马移民局递交永久居留申请，包括护照、无犯罪记录证明、体检报告、房产证等。', 3),
(3, 4, '第四阶段', '1-2个月', '获批与领取居留卡', '移民局审核通过，领取巴拿马永久居留卡（Carné de Residencia Permanente）。5年后可申请入籍。', 4);

-- Seed FAQs
INSERT INTO `faqs` (`project_id`, `question`, `answer`, `is_global`, `sort_order`) VALUES
(1, 'EB-5投资移民的最低投资额是多少？', '根据2022年《EB-5改革与诚信法案》，目标就业区（TEA）的最低投资额为80万美元，非TEA地区为105万美元。', 0, 1),
(1, 'EB-5绿卡多久可以获批？', '目前I-526E审批时间约为12-18个月，加上领事馆面签或身份调整，整体周期约为2-3年。', 0, 2),
(1, 'EB-5投资资金可以返还吗？', '可以。投资期满（通常为5年）后，只要项目创造了10个全职就业岗位，投资款可按项目协议返还。但投资必须"有风险"（at risk），不能有保底承诺。', 0, 3),
(2, '香港CIES最新的投资门槛是多少？', '2024年3月1日起，香港资本投资者入境计划的新投资门槛为3000万港元。', 0, 4),
(2, 'CIES申请人需要满足居住要求吗？', '要申请香港永久性居民身份，需连续在港居住满7年。但续期逗留签证无严格居住要求。', 0, 5),
(3, '巴拿马购房移民的最低投资额是多少？', '巴拿马购房移民要求购买价值30万美元以上的房产，即可申请永久居留权。', 0, 6),
(3, '巴拿马永久居留后可以申请入籍吗？', '获得巴拿马永久居留权5年后，满足居住要求和其他条件后可申请入籍。巴拿马承认双重国籍。', 0, 7),
(NULL, '北极星移民提供哪些服务？', '北极星移民专注于美国EB-5、香港投资移民和巴拿马购房移民三大项目，提供从项目评估、文件准备、申请递交到获批登陆的全流程服务。', 1, 8),
(NULL, '投资移民的成功率如何？', '投资移民的成功率取决于多个因素，包括资金来源的合法性、项目选择的合理性、文件准备的完整性等。北极星移民团队拥有丰富的经验和专业的律师团队，确保每个申请都符合移民局的要求。', 1, 9);

-- Seed cases
INSERT INTO `cases` (`project_id`, `name`, `country_from`, `investment_amount`, `investment_value`, `processing_period`, `description`, `sort_order`) VALUES
(1, '张先生', '中国', '80万美元', 800000.00, '2年3个月', '张先生通过北极星移民办理EB-5投资移民，选择乡村地区项目，享受预留签证快速审理。从递交I-526E到获批移民签证仅用时2年3个月。', 1),
(1, '李女士', '中国', '80万美元', 800000.00, '2年8个月', '李女士为子女教育选择EB-5投资移民，投资TEA地区商业地产项目。获批后子女以本地生身份入读美国公立高中。', 2),
(2, '王先生', '中国', '3000万港元', 3870000.00, '8个月', '王先生通过北极星移民办理香港CIES，投资组合包含港股、债券和指定投资组合。从递交申请到获得正式批准共用时8个月。', 3),
(2, '陈女士', '中国', '3000万港元', 3870000.00, '10个月', '陈女士为扩展海外业务选择香港CIES，投资以金融资产为主。获批后已成功在香港设立公司并开展业务。', 4),
(3, '刘先生', '中国', '30万美元', 300000.00, '4个月', '刘先生选择巴拿马购房移民，在巴拿马城购入一套海景公寓。从购房到获得永久居留卡共用时4个月。', 5),
(3, '赵女士', '中国', '35万美元', 350000.00, '5个月', '赵女士为退休生活选择巴拿马移民，购入一套带花园的别墅。巴拿马气候宜人、生活成本适中，是理想的退休目的地。', 6);

-- Seed home config (hero slides + advantage items)
INSERT INTO `home_configs` (`config_key`, `config_value`) VALUES
('hero_slides', '[
  {"title":"美国EB-5投资移民","desc":"80万美元投资 · 全家绿卡 · 无语言学历要求","project_slug":"eb5","gradient":"linear-gradient(135deg, #1a3a5c 0%, #2d5a8e 50%, #c8963e 100%)","image":""},
  {"title":"香港资本投资者入境计划","desc":"3000万港元投资 · 获香港居留权 · 国际金融中心","project_slug":"cies","gradient":"linear-gradient(135deg, #c8963e 0%, #d4a94e 50%, #1a3a5c 100%)","image":""},
  {"title":"巴拿马购房移民","desc":"30万美元购房 · 3-6个月获批 · 无居住要求","project_slug":"panama","gradient":"linear-gradient(135deg, #1a3a5c 0%, #1a5c3a 50%, #c8963e 100%)","image":""}
]'),
('advantage_items', '[
  {"icon":"shield","icon_type":"lucide","title":"合规保障","description":"所有项目均经律师团队审核，符合移民局最新政策要求"},
  {"icon":"clock","icon_type":"lucide","title":"高效办理","description":"专业文案团队，材料整理准确率高，缩短审核周期"},
  {"icon":"users","icon_type":"lucide","title":"全流程陪伴","description":"从项目选择到获批登陆，全程一对一专属顾问服务"},
  {"icon":"globe","icon_type":"lucide","title":"国际资源","description":"与全球顶级移民律师、区域中心深度合作"}
]');

INSERT INTO `pages` (`title`, `slug`, `content`, `meta_title`, `meta_description`, `template`, `status`, `sort_order`) VALUES
    ('关于我们', 'about', '', '', '', 'default', 'published', 0);
