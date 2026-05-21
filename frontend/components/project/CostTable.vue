<template>
  <div class="cost-table-wrapper">
    <table class="cost-table">
      <thead>
        <tr>
          <th>费用项目</th>
          <th>金额</th>
          <th>说明</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, index) in rows" :key="index">
          <td class="col-item">{{ row.item || '—' }}</td>
          <td class="col-amount">{{ row.amount || '—' }}</td>
          <td class="col-note">{{ row.note || '—' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
interface CostRow {
  item: string;
  amount: string;
  note: string;
}

defineProps<{
  rows: CostRow[];
}>();
</script>

<style scoped>
.cost-table-wrapper {
  overflow-x: auto;
}

.cost-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.cost-table thead {
  background-color: var(--primary);
  color: var(--bg-white);
}

.cost-table th {
  padding: 14px 16px;
  font-weight: 600;
  text-align: left;
  white-space: nowrap;
}

.cost-table th:first-child {
  border-radius: var(--radius-md) 0 0 0;
}

.cost-table th:last-child {
  border-radius: 0 var(--radius-md) 0 0;
}

.cost-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.cost-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.cost-table tbody tr:hover {
  background-color: var(--bg-gray);
}

.col-item {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
}

.col-amount {
  font-weight: 600;
  color: var(--accent-dark);
  white-space: nowrap;
}

.col-note {
  color: var(--text-light);
  font-size: 13px;
}
</style>
