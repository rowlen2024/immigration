<template>
  <div>
    <AdminPageHeader title="项目管理">
      <template #actions>
        <el-button type="primary" @click="openCreate">新建项目</el-button>
      </template>
    </AdminPageHeader>

    <!-- Toolbar -->
    <AdminToolbar>
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
      <el-button :icon="Refresh" circle @click="searchQuery='';statusFilter='';loadList()" :loading="loading" />
    </AdminToolbar>

    <!-- Table -->
    <AdminTableShell :loading="loading">
      <el-table :data="list">
        <el-table-column prop="name" label="项目名称" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="admin-row-title">{{ row.name }}</div>
              <div class="admin-row-meta">{{ row.country }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="investment_amount" label="投资金额" width="130">
          <template #default="{ row }">{{ row.investment_amount || '—' }}</template>
        </el-table-column>
        <el-table-column prop="processing_period" label="办理周期" width="120">
          <template #default="{ row }">{{ row.processing_period || '—' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'draft']">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="updated_at" label="修改时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <AdminRowActions>
              <button class="action-btn" type="button" title="编辑" aria-label="编辑" @click="openEdit(row)" v-html="getIconSvg('pencil', 16)"></button>
              <el-popconfirm
                title="确定删除该项目？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="handleDelete(row.id)"
              >
                <template #reference>
                  <button class="action-btn danger" type="button" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                </template>
              </el-popconfirm>
            </AdminRowActions>
          </template>
        </el-table-column>
      </el-table>
    </AdminTableShell>

    <!-- Empty state -->
    <AdminEmptyState
      v-if="!loading && list.length === 0"
      icon="folder"
      title="暂无项目"
      description="点击上方按钮创建第一个项目"
      action-label="新建项目"
      @action="openCreate"
    />

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
      size="750px"
      destroy-on-close
      @opened="onDialogOpened"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <section class="project-form-section">
              <h3 class="project-form-section-title">基础信息</h3>
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
            </section>
            <section class="project-form-section">
              <h3 class="project-form-section-title">金额与周期</h3>
            <el-row :gutter="12">
              <el-col :span="12">
                <el-form-item label="投资金额" prop="investment_amount">
                  <el-input v-model="form.investment_amount" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="投资价值" prop="investment_value">
                  <el-input-number v-model="form.investment_value" :min="0" :precision="2" class="w-full" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="办理周期" prop="processing_period">
              <el-input v-model="form.processing_period" />
            </el-form-item>
            <el-form-item label="目标人群" prop="target_crowd">
              <el-input v-model="form.target_crowd" />
            </el-form-item>
            </section>
            <section class="project-form-section">
              <h3 class="project-form-section-title">内容文案</h3>
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
            </section>
            <section class="project-form-section">
              <h3 class="project-form-section-title">Hero 展示</h3>
            <el-form-item label="Hero 标题" prop="hero_title">
              <el-input v-model="form.hero_title" />
            </el-form-item>
            <el-form-item label="封面图片">
              <ImageInput v-model="form.cover_image" placeholder="图片 URL 或上传" size-hint="推荐 800×450px (16:9 横向)" context="project" />
            </el-form-item>
            <el-form-item label="Hero 描述" prop="hero_desc">
              <el-input v-model="form.hero_desc" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="Hero 渐变" prop="hero_gradient">
              <el-input v-model="form.hero_gradient" />
            </el-form-item>
            </section>
            <section class="project-form-section">
              <h3 class="project-form-section-title">发布设置</h3>
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
            </section>
          </el-form>
        </el-tab-pane>

        <el-tab-pane v-if="canReadProjectChildren" label="申请条件" name="requirements">
          <AdminSubResourceActions>
            <el-button v-if="canWriteProjectChildren" type="primary" size="small" @click="openSubDialog('requirement')">添加条件</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.requirements" border size="small">
            <el-table-column prop="label" label="条件描述" min-width="180">
              <template #default="{ row: r }">{{ r.label || '—' }}</template>
            </el-table-column>
            <el-table-column label="必需" width="70">
              <template #default="{ row: r }">
                <span :class="['status-pill', r.is_required ? 'published' : 'draft']">
                  {{ r.is_required ? '必需' : '可选' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('requirement', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('requirement', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadProjectChildren" label="费用明细" name="costItems">
          <AdminSubResourceActions>
            <el-button v-if="canWriteProjectChildren" type="primary" size="small" @click="openSubDialog('costItem')">添加费用</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.costItems" border size="small">
            <el-table-column prop="name" label="费用名称" min-width="140">
              <template #default="{ row: r }">{{ r.name || '—' }}</template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row: r }">{{ r.amount || '—' }}</template>
            </el-table-column>
            <el-table-column prop="note" label="说明" min-width="160">
              <template #default="{ row: r }">{{ r.note || '—' }}</template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('costItem', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('costItem', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadProjectChildren" label="申请流程" name="timelinePhases">
          <AdminSubResourceActions>
            <el-button v-if="canWriteProjectChildren" type="primary" size="small" @click="openSubDialog('timelinePhase')">添加步骤</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.timelinePhases" border size="small">
            <el-table-column prop="phase_number" label="步骤号" width="70">
              <template #default="{ row: r }">{{ r.phase_number ?? '—' }}</template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="130">
              <template #default="{ row: r }">{{ r.title || '—' }}</template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="160">
              <template #default="{ row: r }">{{ r.description || '—' }}</template>
            </el-table-column>
            <el-table-column prop="duration" label="周期" width="90">
              <template #default="{ row: r }">{{ r.duration || '—' }}</template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('timelinePhase', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('timelinePhase', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadProjectChildren" label="项目优势" name="advantages">
          <AdminSubResourceActions>
            <el-button v-if="canWriteProjectChildren" type="primary" size="small" @click="openSubDialog('advantage')">添加优势</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.advantages" border size="small">
            <el-table-column label="图标" width="70">
              <template #default="{ row: r }">
                <span
                  v-if="getIconByName(r.icon)"
                  v-html="getIconSvg(r.icon, 18)"
                  class="project-icon-preview"
                ></span>
                <span v-else>{{ r.icon }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="130">
              <template #default="{ row: r }">{{ r.title || '—' }}</template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="160">
              <template #default="{ row: r }">{{ r.description || '—' }}</template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('advantage', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('advantage', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadCases" label="成功案例" name="cases">
          <AdminSubResourceActions>
            <el-button v-if="canWriteCases" type="primary" size="small" @click="openSubDialog('caseItem')">添加案例</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.cases" border size="small" max-height="360">
            <el-table-column prop="name" label="名称" min-width="120">
              <template #default="{ row: r }">{{ r.name || '—' }}</template>
            </el-table-column>
            <el-table-column prop="country_from" label="来源国" width="80">
              <template #default="{ row: r }">{{ r.country_from || '—' }}</template>
            </el-table-column>
            <el-table-column prop="investment_amount" label="投资金额" width="100">
              <template #default="{ row: r }">{{ r.investment_amount || '—' }}</template>
            </el-table-column>
            <el-table-column prop="processing_period" label="处理周期" width="90">
              <template #default="{ row: r }">{{ r.processing_period || '—' }}</template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('caseItem', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('caseItem', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadTestimonials" label="客户评价" name="testimonials">
          <AdminSubResourceActions>
            <el-button v-if="canWriteTestimonials" type="primary" size="small" @click="openSubDialog('testimonial')">添加评价</el-button>
          </AdminSubResourceActions>
          <el-table :data="subData.testimonials" border size="small" max-height="360">
            <el-table-column label="头像" width="64">
              <template #default="{ row }">
                <ResponsiveImage v-if="row.avatar_url" :src="row.avatar_url" variant="thumb" class="thumb-preview" />
                <span v-else class="admin-no-thumb">—</span>
              </template>
            </el-table-column>
            <el-table-column prop="nickname" label="昵称" width="120">
              <template #default="{ row: r }">{{ r.nickname || '—' }}</template>
            </el-table-column>
            <el-table-column label="星级" width="140">
              <template #default="{ row }">
                <el-rate v-model="row.rating" disabled show-score size="small" />
              </template>
            </el-table-column>
            <el-table-column prop="content" label="评价内容" min-width="200" show-overflow-tooltip>
              <template #default="{ row: r }">{{ r.content || '—' }}</template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60">
              <template #default="{ row: r }">{{ r.sort_order ?? '—' }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" title="编辑" aria-label="编辑" @click="openSubDialog('testimonial', r)" v-html="getIconSvg('pencil', 16)"></button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('testimonial', r.id)">
                    <template #reference>
                      <button class="action-btn danger" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadPages" label="资讯" name="news">
          <AdminSubResourceActions>
            <el-button v-if="canWritePages" type="primary" size="small" @click="openNewsDialog">添加资讯</el-button>
          </AdminSubResourceActions>
          <el-table :data="subNews" border size="small" max-height="360">
            <el-table-column label="标题" min-width="180">
              <template #default="{ row }">
                <span>{{ row.title || '—' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <span :class="['status-pill', row.status === 'published' ? 'published' : 'draft']">
                  {{ row.status === 'published' ? '已发布' : '草稿' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80" fixed="right">
              <template #default="{ row }">
                <div class="table-actions">
                  <el-popconfirm title="确定解除关联？" @confirm="removeNewsLink(row.id)">
                    <template #reference>
                      <button class="action-btn danger" title="移除" aria-label="移除" v-html="getIconSvg('x', 16)"></button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane v-if="canReadProjectChildren" label="项目对比" name="compare">
          <AdminSubResourceActions>
            <el-button type="primary" size="small" :disabled="compareConfig.compare_with.length < 2" @click="saveCompareConfig">保存对比配置</el-button>
          </AdminSubResourceActions>
          <el-form label-position="top">
            <el-form-item label="对比项目（至少 2 个）">
              <el-select v-model="compareConfig.compare_with" multiple filterable placeholder="选择对比项目" class="admin-full-width">
                <el-option v-for="p in projectOptions" :key="p.slug" :label="p.name" :value="p.slug" />
              </el-select>
              <div v-if="compareConfig.compare_with.length < 2" class="compare-helper compare-helper--danger">
                当前项目默认在内，请至少追加 1 个其他项目
              </div>
              <div v-else class="compare-helper">
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
        <AdminDrawerFooter
          :loading="saving"
          @cancel="drawerVisible = false"
          @confirm="handleSave"
        />
      </template>
    </el-drawer>

    <!-- Sub-entity dialog -->
    <el-dialog
      v-model="subDialogVisible"
      :title="subDialogTitle"
      :width="subType === 'caseItem' ? '860px' : '500px'"
      :top="subType === 'caseItem' ? '3vh' : '15vh'"
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
          <el-form-item label="内容" prop="content">
            <RichEditor v-model="subForm.content" />
          </el-form-item>
          <el-form-item label="封面图片" prop="photo_url">
            <ImageInput v-model="subForm.photo_url" placeholder="图片 URL 或上传" size-hint="推荐 800×450px (16:9 横向)" context="project" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'advantage'">
          <el-form-item label="图标" prop="icon">
            <IconPicker v-model="subForm.icon" />
          </el-form-item>
          <el-form-item label="标题" prop="title">
            <el-input v-model="subForm.title" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="subForm.description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'testimonial'">
          <el-form-item label="头像" prop="avatar_url">
            <ImageInput v-model="subForm.avatar_url" placeholder="图片 URL 或上传" size-hint="推荐 200×200px (1:1 正方形)" preview-ratio="1 / 1" context="testimonial" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="subForm.nickname" />
          </el-form-item>
          <el-form-item label="星级评分" prop="rating">
            <el-rate v-model="subForm.rating" show-score />
          </el-form-item>
          <el-form-item label="评价内容" prop="content">
            <el-input v-model="subForm.content" type="textarea" :rows="4" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <AdminDrawerFooter
          :loading="subSaving"
          @cancel="subDialogVisible = false"
          @confirm="handleSubSave"
        />
      </template>
    </el-dialog>

    <el-dialog v-model="newsDialogVisible" title="添加资讯" width="500px" destroy-on-close>
      <el-select v-model="newsSelected" multiple filterable placeholder="搜索新闻页面..." class="admin-full-width">
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
import { useNotify } from '~/composables/useNotify';
import ImageInput from '~/components/admin/ImageInput.vue';
import IconPicker from '~/components/admin/IconPicker.vue';
import RichEditor from '~/components/RichEditor.vue';
import { getIconByName, getIconSvg } from '~/composables/lucideIcons';
import { formatDateTime } from '~/utils/date';
import { generateSlugFromText } from '~/utils/slug';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();
const { hasPermission } = usePermissions();

interface Project {
  id: string;
  slug: string;
  name: string;
  country: string;
  flag_emoji: string;
  tagline: string;
  investment_amount: string;
  investment_value: number;
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
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

interface CaseItem {
  id: number;
  name: string;
  country_from: string;
  investment_amount: string;
  processing_period: string;
  content: string;
  photo_url: string;
  sort_order: number;
}

interface NewsItem {
  id: number;
  title: string;
  slug: string;
  status?: string;
  created_at?: string;
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

const defaultForm = () => ({
  id: undefined as string | undefined,
  slug: '',
  name: '',
  country: '',
  flag_emoji: '',
  tagline: '',
  investment_amount: '',
  investment_value: 0,
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
} as Project);

const form = reactive(defaultForm());

const rules: FormRules = {
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  country: [{ required: true, message: '请输入国家', trigger: 'blur' }],
};

type SubType = 'requirement' | 'costItem' | 'timelinePhase' | 'caseItem' | 'advantage' | 'testimonial';
const subTypeLabels: Record<SubType, string> = {
  requirement: '申请条件',
  costItem: '费用明细',
  timelinePhase: '申请流程',
  caseItem: '成功案例',
  advantage: '项目优势',
  testimonial: '客户评价',
};

interface SubState {
  requirements: any[];
  costItems: any[];
  timelinePhases: any[];
  cases: any[];
  advantages: any[];
  testimonials: any[];
}

const subData = reactive<SubState>({
  requirements: [],
  costItems: [],
  timelinePhases: [],
  cases: [],
  advantages: [],
  testimonials: [],
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

const subNews = ref<NewsItem[]>([]);
const newsDialogVisible = ref(false);
const newsOptions = ref<NewsItem[]>([]);
const newsSelected = ref<number[]>([]);

const compareFields = ref<CompareField[]>([]);
const compareConfig = reactive<CompareConfig>({ compare_with: [], compare_fields: [] });
const projectOptions = ref<ProjectOption[]>([]);

const canReadProjectChildren = computed(() => hasPermission('projects:read'));
const canWriteProjectChildren = computed(() => hasPermission('projects:write'));
const canReadCases = computed(() => hasPermission('cases:read'));
const canWriteCases = computed(() => hasPermission('cases:write'));
const canReadTestimonials = computed(() => hasPermission('testimonials:read'));
const canWriteTestimonials = computed(() => hasPermission('testimonials:write'));
const canReadPages = computed(() => hasPermission('pages:read'));
const canWritePages = computed(() => hasPermission('pages:write'));

const defaultSubForm = (type: SubType): Record<string, any> => {
  switch (type) {
    case 'requirement':
      return { label: '', is_required: true, sort_order: 0 };
    case 'costItem':
      return { name: '', amount: '', note: '', sort_order: 0 };
    case 'timelinePhase':
      return { phase_number: 1, title: '', description: '', duration: '', sort_order: 0 };
    case 'caseItem':
      return { name: '', country_from: '', investment_amount: '', processing_period: '', content: '', photo_url: '', sort_order: 0 };
    case 'advantage':
      return { icon: '', icon_type: 'lucide', title: '', description: '', sort_order: 0 };
    case 'testimonial':
      return { avatar_url: '', nickname: '', rating: 5, content: '', sort_order: 0 };
  }
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/projects?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&name=${encodeURIComponent(searchQuery.value)}`;
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
  const { id, created_at, updated_at, deleted_at, ...cleanRow } = row;
  Object.assign(form, cleanRow);
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
    const [reqs, costs, phases, cases, advs, testimonials] = await Promise.all([
      api<any[]>(`/admin/projects/${editingId.value}/requirements`),
      api<any[]>(`/admin/projects/${editingId.value}/cost-items`),
      api<any[]>(`/admin/projects/${editingId.value}/timeline-phases`),
      api<CaseItem[]>(`/admin/projects/${editingId.value}/cases`),
      api<any[]>(`/admin/projects/${editingId.value}/advantages`),
      api<any[]>(`/admin/projects/${editingId.value}/testimonials`),
    ]);
    subData.requirements = reqs ?? [];
    subData.costItems = costs ?? [];
    subData.timelinePhases = phases ?? [];
    subData.cases = cases ?? [];
    subData.advantages = advs ?? [];
    subData.testimonials = testimonials ?? [];
  } catch {
    // best-effort sub-data load
  }
};

const resetSubData = () => {
  subData.requirements = [];
  subData.costItems = [];
  subData.timelinePhases = [];
  subData.cases = [];
  subData.advantages = [];
  subData.testimonials = [];
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
      notify.success('更新成功');
    } else {
      const created = await api<{ id: string }>('/admin/projects', { method: 'POST', body: form });
      if (created?.id) {
        editingId.value = String(created.id);
      }
      notify.success('创建成功');
    }
    drawerVisible.value = false;
    loadList();
  } catch (e) {
    notify.error(e, editingId.value ? '更新项目失败' : '创建项目失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/projects/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    loadList();
  } catch (e) {
    notify.error(e, '删除项目失败');
  }
};

const generateSlug = () => {
  if (!form.name || !form.name.trim()) {
    ElMessage.warning('请先输入项目名称');
    return;
  }
  const slug = generateSlugFromText(form.name);
  if (!slug) {
    ElMessage.warning('未识别到可生成 slug 的有效内容');
    return;
  }
  form.slug = slug;
};

const openSubDialog = (type: SubType, row?: any) => {
  if (!canWriteSubType(type)) {
    ElMessage.warning('无权限操作');
    return;
  }
  subType.value = type;
  subEditingId.value = row?.id ?? null;
  if (row) {
    const { id, created_at, updated_at, deleted_at, project_id, project, ...cleanRow } = row;
    Object.assign(subForm, cleanRow);
  } else {
    Object.assign(subForm, defaultSubForm(type));
    delete subForm.id;
  }
  if (type === 'requirement') {
    subForm.is_required = subForm.is_required === true || subForm.is_required === 1;
  }
  subDialogVisible.value = true;
};

const canWriteSubType = (type: SubType) => {
  if (type === 'caseItem') return canWriteCases.value;
  if (type === 'testimonial') return canWriteTestimonials.value;
  return canWriteProjectChildren.value;
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
    else if (subType.value === 'advantage') endpoint += '/advantages';
    else if (subType.value === 'testimonial') endpoint += '/testimonials';
    else endpoint += '/timeline-phases';

    if (subEditingId.value) {
      endpoint += `/${subEditingId.value}`;
      await api(endpoint, { method: 'PUT', body: subForm });
    } else {
      await api(endpoint, { method: 'POST', body: subForm });
    }
    notify.success(subEditingId.value ? '更新成功' : '添加成功');
    subDialogVisible.value = false;
    loadSubData();
  } catch (e) {
    notify.error(e, '保存失败');
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
    else if (type === 'advantage') endpoint += `/advantages/${id}`;
    else if (type === 'testimonial') endpoint += `/testimonials/${id}`;
    else endpoint += `/timeline-phases/${id}`;
    await api(endpoint, { method: 'DELETE' });
    notify.success('已删除');
    loadSubData();
  } catch (e) {
    notify.error(e, '删除失败');
  }
};

const loadNews = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    const data = await api<NewsItem[]>(`/admin/projects/${editingId.value}/news`);
    subNews.value = data ?? [];
  } catch { subNews.value = []; }
};

const openNewsDialog = async () => {
  if (!canWritePages.value) {
    ElMessage.warning('无权限操作');
    return;
  }
  try {
    const api = useApi();
    const data = await api<{ items: NewsItem[] }>('/admin/pages/options?page=1&per_page=500&page_type=news&status=published');
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
    notify.success('已关联');
    newsDialogVisible.value = false;
    loadNews();
  } catch (e) { notify.error(e, '添加资讯失败'); }
};

const removeNewsLink = async (pageId: number) => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/news/${pageId}`, { method: 'DELETE' });
    notify.success('已解除关联');
    loadNews();
  } catch (e) { notify.error(e, '解除关联失败'); }
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
    const data = await api<{ items: ProjectOption[] }>('/admin/projects/options?page=1&per_page=500');
    projectOptions.value = data.items ?? [];
  } catch { projectOptions.value = []; }
};

const saveCompareConfig = async () => {
  if (!editingId.value) return;
  if (compareConfig.compare_with.length < 2) {
    ElMessage.warning('请至少选择 2 个对比项目');
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
    notify.success('项目对比配置已保存');
  } catch (e) { notify.error(e, '保存对比配置失败'); }
};


onMounted(() => {
  loadList();
});
</script>

<style scoped>
.project-form-section {
  padding: 16px;
  margin-bottom: 14px;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.project-form-section-title {
  margin: 0 0 14px;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.thumb-preview {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 50%;
}

.project-icon-preview {
  display: inline-flex;
  align-items: center;
  color: var(--color-accent);
}

.compare-helper {
  margin-top: 4px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.compare-helper--danger {
  color: var(--color-danger);
}
</style>
