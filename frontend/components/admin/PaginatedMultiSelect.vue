<template>
  <el-select
    v-model="selected"
    multiple
    filterable
    remote
    :remote-method="searchOptions"
    :loading="loading"
    :placeholder="placeholder"
    :disabled="disabled"
    @visible-change="onVisibleChange"
  >
    <el-option
      v-for="option in options"
      :key="option.value"
      :label="option.label"
      :value="option.value"
    />
    <template #empty>
      <div class="paginated-select-empty">
        <template v-if="loading">加载中...</template>
        <template v-else-if="loadError">
          加载失败
          <el-button link type="primary" @click.stop="loadCurrentPage">重试</el-button>
        </template>
        <template v-else>{{ emptyText }}</template>
      </div>
    </template>
    <template #footer>
      <el-pagination
        v-if="!loadError"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        :pager-count="5"
        :disabled="loading"
        size="small"
        layout="total, prev, pager, next"
        @update:current-page="changePage"
      />
    </template>
  </el-select>
</template>

<script setup lang="ts" generic="T extends PaginatedSelectValue">
import type {
  PaginatedSelectOption,
  PaginatedSelectRequest,
  PaginatedSelectResult,
  PaginatedSelectValue,
} from '~/types/paginated-select';

/** 通用分页多选配置，业务接口及筛选规则由加载函数负责。 */
const props = withDefaults(defineProps<{
  /** 按当前查询返回一页选项，不得在加载函数中循环查询全量数据。 */
  loadOptions: (request: PaginatedSelectRequest) => Promise<PaginatedSelectResult<T>>;
  /** 每页选项数量，应符合所调用接口的分页约定。 */
  pageSize?: number;
  /** 未输入关键词时的提示文本。 */
  placeholder?: string;
  /** 当前检索没有选项时的提示文本。 */
  emptyText?: string;
  /** 是否禁止修改选择。 */
  disabled?: boolean;
}>(), {
  pageSize: 20,
  placeholder: '输入关键词搜索...',
  emptyText: '未找到匹配项',
  disabled: false,
});

/** 已选择的原始标识，搜索和翻页时保持不变。 */
const selected = defineModel<T[]>({ required: true });
/** 当前页选项，跨页已选项的标题由 Element Plus 的选项缓存保留。 */
const options = ref<PaginatedSelectOption<T>[]>([]);
/** 从 1 开始的当前页码。 */
const page = ref(1);
/** 当前检索条件下的选项总数。 */
const total = ref(0);
/** 原样传递给业务加载函数的关键词。 */
const keyword = ref('');
/** 当前是否正在请求选项。 */
const loading = ref(false);
/** 当前请求是否失败，用于显示重试入口。 */
const loadError = ref(false);
/** 是否已成功加载过选项，避免展开下拉框时重复加载同一页。 */
const loaded = ref(false);
/** 最近发起的请求序号，用于忽略搜索或销毁前的过期响应。 */
let requestId = 0;

/** 请求当前页，仅接受最新请求的结果，加载失败时保留选择并提供重试。 */
const loadCurrentPage = async () => {
  /** 当前请求的快照序号。 */
  const currentRequestId = ++requestId;
  loading.value = true;
  loadError.value = false;
  options.value = [];
  try {
    /** 业务加载函数负责接口信封解包和选项字段映射。 */
    const result = await props.loadOptions({
      page: page.value,
      perPage: props.pageSize,
      keyword: keyword.value,
    });
    // 快速输入和切换分页会产生并行请求，旧响应不能覆盖较新的检索结果。
    if (currentRequestId !== requestId) return;
    options.value = result.items;
    total.value = result.total;
    loaded.value = true;
  } catch {
    if (currentRequestId !== requestId) return;
    loadError.value = true;
  } finally {
    if (currentRequestId === requestId) loading.value = false;
  }
};

/** 首次展开才加载选项，重复展开保留当前页。 */
const onVisibleChange = (visible: boolean) => {
  if (visible && !loaded.value && !loading.value) loadCurrentPage();
};

/** 复用选择器的输入防抖，关键词改变后从第一页检索并保留已选项。 */
const searchOptions = (value: string) => {
  if (value === keyword.value) return;
  keyword.value = value;
  page.value = 1;
  total.value = 0;
  loadCurrentPage();
};

/** 翻页时只加载目标页，已选标识由父页面继续持有。 */
const changePage = (value: number) => {
  page.value = value;
  loadCurrentPage();
};

onBeforeUnmount(() => {
  // 弹窗销毁后返回的请求不能再更新选项或错误状态。
  requestId++;
});
</script>

<style scoped>
.paginated-select-empty {
  padding: 16px;
  color: var(--color-text-muted);
  text-align: center;
}
</style>
