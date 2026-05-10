<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">项目管理</h2>
      <el-button type="primary" @click="openCreate">新建项目</el-button>
    </div>

    <!-- Toolbar -->
    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索项目名称..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-select
        v-model="statusFilter"
        placeholder="状态筛选"
        clearable
        class="admin-filter-select"
        @change="loadList"
      >
        <el-option label="全部" value="" />
        <el-option label="已发布" value="1" />
        <el-option label="草稿" value="0" />
      </el-select>
      <el-button :icon="Refresh" circle @click="loadList" :loading="loading" style="margin-left:auto;" />
    </div>

    <!-- Table -->
    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="name" label="项目名称" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="row-title">{{ row.name }}</div>
              <div class="row-meta">{{ row.country }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="investment_amount" label="投资金额" width="130" />
        <el-table-column prop="processing_period" label="办理周期" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'draft']">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm
                title="确定删除该项目？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="handleDelete(row.id)"
              >
                <template #reference>
                  <button class="action-btn danger">删除</button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Empty state -->
    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('folder', 48)"></div>
      <div class="empty-title">暂无项目</div>
      <div class="empty-desc">点击上方按钮创建第一个项目</div>
      <el-button type="primary" @click="openCreate">新建项目</el-button>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadList"
      />
    </div>

    <!-- Drawer -->
    <el-drawer
      v-model="drawerVisible"
      :title="editingId ? '编辑项目' : '新建项目'"
      size="700px"
      destroy-on-close
      @opened="onDialogOpened"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="标识(slug)" prop="slug">
                  <el-input v-model="form.slug">
                    <template #suffix>
                      <el-button
                        link
                        type="primary"
                        size="small"
                        :disabled="!form.name || !form.name.trim()"
                        @click="generateSlug"
                      >自动生成</el-button>
                    </template>
                  </el-input>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="项目名称" prop="name">
                  <el-input v-model="form.name" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="国家" prop="country">
                  <el-input v-model="form.country" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="国旗图标" prop="flag_emoji">
                  <el-input v-model="form.flag_emoji" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="标语" prop="tagline">
              <el-input v-model="form.tagline" />
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="投资金额" prop="investment_amount">
                  <el-input v-model="form.investment_amount" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="投资价值" prop="investment_value">
                  <el-input v-model="form.investment_value" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="办理周期" prop="processing_period">
              <el-input v-model="form.processing_period" />
            </el-form-item>
            <el-form-item label="目标人群" prop="target_crowd">
              <el-input v-model="form.target_crowd" />
            </el-form-item>
            <el-form-item label="概览标题" prop="overview_title">
              <el-input v-model="form.overview_title" />
            </el-form-item>
            <el-form-item label="概览内容" prop="overview_text">
              <el-input v-model="form.overview_text" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item label="政策标题" prop="policy_title">
              <el-input v-model="form.policy_title" />
            </el-form-item>
            <el-form-item label="政策内容" prop="policy_text">
              <el-input v-model="form.policy_text" type="textarea" :rows="3" />
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="费用总计" prop="costs_total">
                  <el-input v-model="form.costs_total" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="费用备注" prop="costs_note">
                  <el-input v-model="form.costs_note" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="CTA 文字" prop="cta_text">
              <el-input v-model="form.cta_text" />
            </el-form-item>
            <el-form-item label="Hero 标题" prop="hero_title">
              <el-input v-model="form.hero_title" />
            </el-form-item>
            <el-form-item label="封面图片">
              <ImageInput v-model="form.cover_image" placeholder="图片 URL 或上传" />
            </el-form-item>
            <el-form-item label="Hero 描述" prop="hero_desc">
              <el-input v-model="form.hero_desc" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="Hero 渐变" prop="hero_gradient">
              <el-input v-model="form.hero_gradient" />
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="排序" prop="sort_order">
                  <el-input-number v-model="form.sort_order" :min="0" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="状态" prop="status">
                  <el-select v-model="form.status">
                    <el-option label="草稿" :value="0" />
                    <el-option label="已发布" :value="1" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="申请条件" name="requirements">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('requirement')">添加条件</el-button>
          </div>
          <el-table :data="subData.requirements" border size="small">
            <el-table-column prop="label" label="条件描述" min-width="180" />
            <el-table-column label="必需" width="70">
              <template #default="{ row: r }">
                <span :class="['status-pill', r.is_required ? 'published' : 'draft']">
                  {{ r.is_required ? '必需' : '可选' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('requirement', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('requirement', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="费用明细" name="costItems">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('costItem')">添加费用</el-button>
          </div>
          <el-table :data="subData.costItems" border size="small">
            <el-table-column prop="name" label="费用名称" min-width="140" />
            <el-table-column prop="amount" label="金额" width="100" />
            <el-table-column prop="note" label="说明" min-width="160" />
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('costItem', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('costItem', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="申请流程" name="timelinePhases">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('timelinePhase')">添加步骤</el-button>
          </div>
          <el-table :data="subData.timelinePhases" border size="small">
            <el-table-column prop="phase_number" label="步骤号" width="70" />
            <el-table-column prop="title" label="标题" min-width="130" />
            <el-table-column prop="description" label="描述" min-width="160" />
            <el-table-column prop="duration" label="周期" width="90" />
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('timelinePhase', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('timelinePhase', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="成功案例" name="cases">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('caseItem')">添加案例</el-button>
          </div>
          <el-table :data="subData.cases" border size="small" max-height="360">
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column prop="country_from" label="来源国" width="80" />
            <el-table-column prop="investment_amount" label="投资金额" width="100" />
            <el-table-column prop="processing_period" label="处理周期" width="90" />
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('caseItem', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('caseItem', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="资讯" name="news">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openNewsDialog">添加资讯</el-button>
          </div>
          <el-table :data="subNews" border size="small" max-height="360">
            <el-table-column label="标题" min-width="180">
              <template #default="{ row }">
                <span>{{ row.page?.title || '(已删除)' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <span :class="['status-pill', row.page?.status === 'published' ? 'published' : 'draft']">
                  {{ row.page?.status === 'published' ? '已发布' : '草稿' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-popconfirm title="确定解除关联？" @confirm="removeNewsLink(row.page_id)">
                    <template #reference>
                      <button class="action-btn danger">移除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="项目对比" name="compare">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" :disabled="compareConfig.compare_with.length < 2" @click="saveCompareConfig">保存对比配置</el-button>
          </div>
          <el-form label-position="top">
            <el-form-item label="对比项目（至少 2 个）">
              <el-select v-model="compareConfig.compare_with" multiple filterable placeholder="选择对比项目" style="width: 100%">
                <el-option v-for="p in projectOptions" :key="p.slug" :label="p.name" :value="p.slug" />
              </el-select>
              <div v-if="compareConfig.compare_with.length < 2" style="font-size: 12px; color: var(--el-color-danger); margin-top: 4px">
                当前项目默认在内，请至少追加 1 个其他项目
              </div>
              <div v-else style="font-size: 12px; color: var(--el-text-color-secondary); margin-top: 4px">
                已选择 {{ compareConfig.compare_with.length }} 个项目
              </div>
            </el-form-item>
            <el-form-item label="对比属性">
              <el-checkbox-group v-model="compareConfig.compare_fields">
                <el-checkbox v-for="f in compareFields" :key="f.key" :label="f.key" :value="f.key">
                  {{ f.label }}
                </el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>

    <!-- Sub-entity dialog -->
    <el-dialog
      v-model="subDialogVisible"
      :title="subDialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form ref="subFormRef" :model="subForm" label-position="top">
        <template v-if="subType === 'requirement'">
          <el-form-item label="条件描述" prop="label">
            <el-input v-model="subForm.label" />
          </el-form-item>
          <el-form-item label="是否必需" prop="is_required">
            <el-switch v-model="subForm.is_required" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'costItem'">
          <el-form-item label="费用名称" prop="name">
            <el-input v-model="subForm.name" />
          </el-form-item>
          <el-form-item label="金额" prop="amount">
            <el-input v-model="subForm.amount" />
          </el-form-item>
          <el-form-item label="说明" prop="note">
            <el-input v-model="subForm.note" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'timelinePhase'">
          <el-form-item label="步骤号" prop="phase_number">
            <el-input-number v-model="subForm.phase_number" :min="1" />
          </el-form-item>
          <el-form-item label="标题" prop="title">
            <el-input v-model="subForm.title" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="subForm.description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="周期" prop="duration">
            <el-input v-model="subForm.duration" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'caseItem'">
          <el-form-item label="名称" prop="name">
            <el-input v-model="subForm.name" />
          </el-form-item>
          <el-form-item label="来源国" prop="country_from">
            <el-input v-model="subForm.country_from" />
          </el-form-item>
          <el-form-item label="投资金额" prop="investment_amount">
            <el-input v-model="subForm.investment_amount" />
          </el-form-item>
          <el-form-item label="处理周期" prop="processing_period">
            <el-input v-model="subForm.processing_period" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="subForm.description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="照片URL" prop="photo_url">
            <el-input v-model="subForm.photo_url" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="subDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="subSaving" @click="handleSubSave">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="newsDialogVisible" title="添加资讯" width="500px" destroy-on-close>
      <el-select v-model="newsSelected" multiple filterable placeholder="搜索新闻页面..." style="width: 100%">
        <el-option v-for="n in newsOptions" :key="n.id" :label="n.title" :value="n.id" />
      </el-select>
      <template #footer>
        <el-button @click="newsDialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="newsSelected.length === 0" @click="addNewsLinks">确认添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Search, Refresh } from '@element-plus/icons-vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import ImageInput from '~/components/admin/ImageInput.vue';
import { getIconSvg } from '~/composables/lucideIcons';
import { pinyin } from 'pinyin-pro';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Project {
  id: string;
  slug: string;
  name: string;
  country: string;
  flag_emoji: string;
  tagline: string;
  investment_amount: string;
  investment_value: string;
  processing_period: string;
  target_crowd: string;
  overview_title: string;
  overview_text: string;
  policy_title: string;
  policy_text: string;
  costs_total: string;
  costs_note: string;
  cta_text: string;
  hero_title: string;
  hero_desc: string;
  hero_gradient: string;
  cover_image: string;
  sort_order: number;
  status: number;
}

interface CaseItem {
  id: number;
  name: string;
  country_from: string;
  investment_amount: string;
  processing_period: string;
  description: string;
  photo_url: string;
  sort_order: number;
}

interface NewsItem {
  id: number;
  title: string;
  slug: string;
  status: string;
  created_at: string;
}

interface NewsLink {
  page_id: number;
  page?: NewsItem;
  created_at: string;
}

interface CompareField {
  key: string;
  label: string;
  from: string;
}

interface CompareConfig {
  compare_with: string[];
  compare_fields: string[];
}

interface ProjectOption {
  slug: string;
  name: string;
}

const list = ref<Project[]>([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const drawerVisible = ref(false);
const editingId = ref<string | null>(null);
const formRef = ref<FormInstance>();
const activeTab = ref('basic');

const searchQuery = ref('');
const statusFilter = ref('');

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const defaultForm = (): Partial<Project> => ({
  slug: '',
  name: '',
  country: '',
  flag_emoji: '',
  tagline: '',
  investment_amount: '',
  investment_value: '',
  processing_period: '',
  target_crowd: '',
  overview_title: '',
  overview_text: '',
  policy_title: '',
  policy_text: '',
  costs_total: '',
  costs_note: '',
  cta_text: '',
  hero_title: '',
  hero_desc: '',
  hero_gradient: '',
  cover_image: '',
  sort_order: 0,
  status: 0,
});

const form = reactive<Partial<Project>>(defaultForm());

const rules: FormRules = {
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  country: [{ required: true, message: '请输入国家', trigger: 'blur' }],
};

type SubType = 'requirement' | 'costItem' | 'timelinePhase' | 'caseItem';
const subTypeLabels: Record<SubType, string> = {
  requirement: '申请条件',
  costItem: '费用明细',
  timelinePhase: '申请流程',
  caseItem: '成功案例',
};

interface SubState {
  requirements: any[];
  costItems: any[];
  timelinePhases: any[];
  cases: any[];
}

const subData = reactive<SubState>({
  requirements: [],
  costItems: [],
  timelinePhases: [],
  cases: [],
});

const subDialogVisible = ref(false);
const subSaving = ref(false);
const subType = ref<SubType>('requirement');
const subEditingId = ref<number | null>(null);
const subFormRef = ref<FormInstance>();
const subForm = reactive<Record<string, any>>({});
const subDialogTitle = computed(() => {
  const prefix = subEditingId.value ? '编辑' : '新增';
  return `${prefix}${subTypeLabels[subType.value]}`;
});

const subNews = ref<NewsLink[]>([]);
const newsDialogVisible = ref(false);
const newsOptions = ref<NewsItem[]>([]);
const newsSelected = ref<number[]>([]);

const compareFields = ref<CompareField[]>([]);
const compareConfig = reactive<CompareConfig>({ compare_with: [], compare_fields: [] });
const projectOptions = ref<ProjectOption[]>([]);

const defaultSubForm = (type: SubType): Record<string, any> => {
  switch (type) {
    case 'requirement':
      return { label: '', is_required: true, sort_order: 0 };
    case 'costItem':
      return { name: '', amount: '', note: '', sort_order: 0 };
    case 'timelinePhase':
      return { phase_number: 1, title: '', description: '', duration: '', sort_order: 0 };
    case 'caseItem':
      return { name: '', country_from: '', investment_amount: '', processing_period: '', description: '', photo_url: '', sort_order: 0 };
  }
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/projects?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (statusFilter.value) url += `&status=${statusFilter.value}`;
    const data = await api<{ items: Project[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载项目列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  resetSubData();
  activeTab.value = 'basic';
  drawerVisible.value = true;
};

const openEdit = (row: Project) => {
  editingId.value = row.id;
  Object.assign(form, row);
  resetSubData();
  activeTab.value = 'basic';
  drawerVisible.value = true;
};

const onDialogOpened = () => {
  if (editingId.value) {
    loadSubData();
    loadNews();
    loadCompareConfig();
    loadProjectOptions();
  }
};

const loadSubData = async () => {
  if (!editingId.value) return;
  const api = useApi();
  try {
    const [reqs, costs, phases, cases] = await Promise.all([
      api<any[]>(`/admin/projects/${editingId.value}/requirements`),
      api<any[]>(`/admin/projects/${editingId.value}/cost-items`),
      api<any[]>(`/admin/projects/${editingId.value}/timeline-phases`),
      api<CaseItem[]>(`/admin/projects/${editingId.value}/cases`),
    ]);
    subData.requirements = reqs ?? [];
    subData.costItems = costs ?? [];
    subData.timelinePhases = phases ?? [];
    subData.cases = cases ?? [];
  } catch {
    // best-effort sub-data load
  }
};

const resetSubData = () => {
  subData.requirements = [];
  subData.costItems = [];
  subData.timelinePhases = [];
  subData.cases = [];
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/projects/${editingId.value}`, {
        method: 'PUT',
        body: form,
      });
    } else {
      const created = await api<{ id: string }>('/admin/projects', { method: 'POST', body: form });
      if (created?.id) {
        editingId.value = String(created.id);
      }
    }
    drawerVisible.value = false;
    loadList();
  } catch {
    ElMessage.error(editingId.value ? '更新项目失败' : '创建项目失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/projects/${id}`, { method: 'DELETE' });
    loadList();
  } catch {
    ElMessage.error('删除项目失败');
  }
};

const generateSlug = () => {
  if (!form.name || !form.name.trim()) {
    ElMessage.warning('请先输入项目名称');
    return;
  }
  const arr = pinyin(form.name, { toneType: 'none', type: 'array' });
  const nonEmpty = arr.filter((s: string) => s.trim() !== '');
  if (nonEmpty.length === 0) {
    ElMessage.warning('未识别到可生成拼音的文字');
    return;
  }
  form.slug = nonEmpty.join('-').toLowerCase();
};

const openSubDialog = (type: SubType, row?: any) => {
  subType.value = type;
  subEditingId.value = row?.id ?? null;
  Object.assign(subForm, row ? { ...row } : defaultSubForm(type));
  if (type === 'requirement') {
    subForm.is_required = subForm.is_required === true || subForm.is_required === 1;
  }
  subDialogVisible.value = true;
};

const handleSubSave = async () => {
  if (!editingId.value) {
    ElMessage.warning('请先保存项目基本信息');
    return;
  }
  subSaving.value = true;
  try {
    const api = useApi();
    let endpoint = `/admin/projects/${editingId.value}`;
    if (subType.value === 'requirement') endpoint += '/requirements';
    else if (subType.value === 'costItem') endpoint += '/cost-items';
    else if (subType.value === 'caseItem') endpoint += '/cases';
    else endpoint += '/timeline-phases';

    if (subEditingId.value) {
      endpoint += `/${subEditingId.value}`;
      await api(endpoint, { method: 'PUT', body: subForm });
    } else {
      await api(endpoint, { method: 'POST', body: subForm });
    }
    subDialogVisible.value = false;
    loadSubData();
  } catch {
    ElMessage.error('保存失败');
  } finally {
    subSaving.value = false;
  }
};

const deleteSubItem = async (type: SubType, id: number) => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    let endpoint = `/admin/projects/${editingId.value}`;
    if (type === 'requirement') endpoint += `/requirements/${id}`;
    else if (type === 'costItem') endpoint += `/cost-items/${id}`;
    else if (type === 'caseItem') endpoint += `/cases/${id}`;
    else endpoint += `/timeline-phases/${id}`;
    await api(endpoint, { method: 'DELETE' });
    loadSubData();
  } catch {
    ElMessage.error('删除失败');
  }
};

const loadNews = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    const data = await api<NewsLink[]>(`/admin/projects/${editingId.value}/news`);
    subNews.value = data ?? [];
  } catch { subNews.value = []; }
};

const openNewsDialog = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: NewsItem[] }>('/admin/pages?page_type=news&status=published&all=true');
    newsOptions.value = data.items ?? [];
  } catch { newsOptions.value = []; }
  newsSelected.value = [];
  newsDialogVisible.value = true;
};

const addNewsLinks = async () => {
  if (!editingId.value || newsSelected.value.length === 0) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/news`, { method: 'POST', body: { page_ids: newsSelected.value } });
    newsDialogVisible.value = false;
    loadNews();
  } catch { ElMessage.error('添加资讯失败'); }
};

const removeNewsLink = async (pageId: number) => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/news/${pageId}`, { method: 'DELETE' });
    loadNews();
  } catch { ElMessage.error('解除关联失败'); }
};

const loadCompareConfig = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    const [fields, cfg] = await Promise.all([
      api<CompareField[]>('/admin/compare-fields'),
      api<CompareConfig>(`/admin/projects/${editingId.value}/compare-config`),
    ]);
    compareFields.value = fields ?? [];
    if (cfg) {
      compareConfig.compare_with = cfg.compare_with ?? [];
      compareConfig.compare_fields = cfg.compare_fields ?? [];
    } else {
      compareConfig.compare_with = [];
      compareConfig.compare_fields = [];
    }
    // Auto-include current project
    const currentSlug = form.slug || '';
    if (currentSlug && !compareConfig.compare_with.includes(currentSlug)) {
      compareConfig.compare_with.unshift(currentSlug);
    }
  } catch {}
};

const loadProjectOptions = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: ProjectOption[] }>('/admin/projects?all=true');
    projectOptions.value = data.items ?? [];
  } catch { projectOptions.value = []; }
};

const saveCompareConfig = async () => {
  if (!editingId.value) return;
  if (compareConfig.compare_with.length < 2) {
    ElMessage({ message: '请至少选择 2 个对比项目', type: 'warning', zIndex: 9999 });
    return;
  }
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/compare-config`, {
      method: 'PUT',
      body: {
        project_id: Number(editingId.value),
        compare_with: compareConfig.compare_with,
        compare_fields: compareConfig.compare_fields,
      },
    });
    ElMessage.success('项目对比配置已保存');
  } catch { ElMessage.error('保存对比配置失败'); }
};


onMounted(() => {
  loadList();
});
</script>

<style scoped>
.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.row-meta {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 2px;
}
</style>
