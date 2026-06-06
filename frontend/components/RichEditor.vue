<template>
  <div class="rich-editor" :class="{ 'source-mode': sourceMode }">
    <!-- Toolbar -->
    <div v-if="editor && !sourceMode" class="editor-toolbar">
      <button
        v-for="(btn, i) in toolbarButtons"
        :key="i"
        type="button"
        class="toolbar-btn"
        :class="{ active: btn.isActive?.(editor) }"
        :title="btn.title"
        @click="btn.action(editor, $event)"
        v-html="btn.icon"
      />
      <span class="toolbar-divider" />
      <button type="button" class="toolbar-btn" title="源码" @click="toggleSource">
        &lt;/&gt;
      </button>
    </div>

    <!-- Editor content -->
    <EditorContent v-show="!sourceMode" :editor="editor" class="editor-content" />

    <!-- Source mode textarea -->
    <textarea
      v-show="sourceMode"
      v-model="sourceHtml"
      class="editor-source"
      rows="16"
    />

    <!-- Hidden file input for image upload -->
    <input
      ref="imageInput"
      type="file"
      accept=".jpg,.jpeg,.png,.webp,.gif,.svg"
      hidden
      @change="handleImageUpload"
    />

    <!-- Table size picker -->
    <div v-if="tablePicker.visible" class="table-picker-popup" :style="{ left: tablePicker.x + 'px', top: tablePicker.y + 'px' }">
      <div class="table-picker-grid" @mouseleave="tablePicker.visible = false">
        <template v-for="r in 6" :key="r">
          <div
            v-for="c in 6"
            :key="c"
            class="grid-cell"
            :class="{ selected: r <= tablePicker.rows && c <= tablePicker.cols }"
            @mouseenter="tablePicker.rows = r; tablePicker.cols = c"
            @click="insertTable(r, c)"
          />
        </template>
      </div>
      <div class="table-picker-label">{{ tablePicker.rows }} × {{ tablePicker.cols }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { EditorContent, useEditor } from '@tiptap/vue-3';
import { type Editor } from '@tiptap/core';
import { StarterKit } from '@tiptap/starter-kit';
import { Table } from '@tiptap/extension-table';
import { TableRow } from '@tiptap/extension-table-row';
import { TableCell } from '@tiptap/extension-table-cell';
import { TableHeader } from '@tiptap/extension-table-header';
import { Image } from '@tiptap/extension-image';
import { TextStyle } from '@tiptap/extension-text-style';
import { Color } from '@tiptap/extension-color';
import { Highlight } from '@tiptap/extension-highlight';
import { TextAlign } from '@tiptap/extension-text-align';
import { Youtube } from '@tiptap/extension-youtube';
import { Underline } from '@tiptap/extension-underline';
import { CodeBlockLowlight } from '@tiptap/extension-code-block-lowlight';
import { common, createLowlight } from 'lowlight';

const lowlight = createLowlight(common);

const props = defineProps<{ modelValue: string }>();
const emit = defineEmits<{ 'update:modelValue': [value: string] }>();

const imageInput = ref<HTMLInputElement>();
const sourceMode = ref(false);
const sourceHtml = ref('');

const tablePicker = reactive({ visible: false, x: 0, y: 0, rows: 3, cols: 3 });

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({ codeBlock: false }),
    Underline,
    Table.configure({ resizable: true }),
    TableRow,
    TableCell,
    TableHeader,
    Image.configure({ allowBase64: false }),
    TextStyle,
    Color,
    Highlight.configure({ multicolor: true }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Youtube.configure({ controls: true }),
    CodeBlockLowlight.configure({ lowlight }),
  ],
  onUpdate: ({ editor: ed }) => {
    emit('update:modelValue', ed.getHTML());
  },
});

watch(
  () => props.modelValue,
  (val) => {
    if (editor.value && editor.value.getHTML() !== val) {
      editor.value.commands.setContent(val, { emitUpdate: false });
    }
  }
);

const getToken = () => {
  if (import.meta.client) {
    return localStorage.getItem('token') || '';
  }
  return '';
};

const handleImageUpload = async (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  const form = new FormData();
  form.append('file', file);
  try {
    const res = await $fetch<{ code: number; data: { url: string } }>('/api/v1/admin/media/upload', {
      method: 'POST',
      headers: { Authorization: `Bearer ${getToken()}` },
      body: form,
    });
    if (res?.data?.url) {
      editor.value?.chain().focus().setImage({ src: res.data.url }).run();
    }
  } catch {
    // ignore
  }
  if (imageInput.value) imageInput.value.value = '';
};

const openImageUpload = () => {
  imageInput.value?.click();
};

const openLinkDialog = () => {
  const url = window.prompt('链接地址 (https://...)');
  if (url) {
    editor.value?.chain().focus().setLink({ href: url }).run();
  }
};

const openYoutubeDialog = () => {
  const url = window.prompt('YouTube 视频链接');
  if (url) {
    editor.value?.chain().focus().setYoutubeVideo({ src: url }).run();
  }
};

const openTablePicker = (e: MouseEvent) => {
  const btn = e.currentTarget as HTMLElement;
  const rect = btn.getBoundingClientRect();
  tablePicker.x = rect.left;
  tablePicker.y = rect.bottom + 4;
  tablePicker.visible = true;
  tablePicker.rows = 3;
  tablePicker.cols = 3;
};

const insertTable = (rows: number, cols: number) => {
  editor.value?.chain().focus().insertTable({ rows, cols, withHeaderRow: true }).run();
  tablePicker.visible = false;
};

const toggleSource = () => {
  if (sourceMode.value) {
    editor.value?.commands.setContent(sourceHtml.value, { emitUpdate: false });
    sourceMode.value = false;
  } else {
    sourceHtml.value = editor.value?.getHTML() ?? '';
    sourceMode.value = true;
  }
};

type EditorInstance = Editor | undefined;

type ToolbarButton = {
  title: string;
  icon: string;
  action: (editor: EditorInstance, event: MouseEvent) => void;
  isActive?: (editor: EditorInstance) => boolean;
};

const toolbarButtons: ToolbarButton[] = [
  {
    title: '加粗',
    icon: '<b>B</b>',
    action: (e) => e?.chain().focus().toggleBold().run(),
    isActive: (e) => e?.isActive('bold') ?? false,
  },
  {
    title: '斜体',
    icon: '<i>I</i>',
    action: (e) => e?.chain().focus().toggleItalic().run(),
    isActive: (e) => e?.isActive('italic') ?? false,
  },
  {
    title: '下划线',
    icon: '<u>U</u>',
    action: (e) => e?.chain().focus().toggleUnderline().run(),
    isActive: (e) => e?.isActive('underline') ?? false,
  },
  {
    title: '标题1',
    icon: 'H1',
    action: (e) => e?.chain().focus().toggleHeading({ level: 1 }).run(),
    isActive: (e) => e?.isActive('heading', { level: 1 }) ?? false,
  },
  {
    title: '标题2',
    icon: 'H2',
    action: (e) => e?.chain().focus().toggleHeading({ level: 2 }).run(),
    isActive: (e) => e?.isActive('heading', { level: 2 }) ?? false,
  },
  {
    title: '标题3',
    icon: 'H3',
    action: (e) => e?.chain().focus().toggleHeading({ level: 3 }).run(),
    isActive: (e) => e?.isActive('heading', { level: 3 }) ?? false,
  },
  {
    title: '引用',
    icon: '❝',
    action: (e) => e?.chain().focus().toggleBlockquote().run(),
    isActive: (e) => e?.isActive('blockquote') ?? false,
  },
  {
    title: '无序列表',
    icon: '•',
    action: (e) => e?.chain().focus().toggleBulletList().run(),
    isActive: (e) => e?.isActive('bulletList') ?? false,
  },
  {
    title: '有序列表',
    icon: '1.',
    action: (e) => e?.chain().focus().toggleOrderedList().run(),
    isActive: (e) => e?.isActive('orderedList') ?? false,
  },
  {
    title: '代码块',
    icon: '&lt;&gt;',
    action: (e) => e?.chain().focus().toggleCodeBlock().run(),
    isActive: (e) => e?.isActive('codeBlock') ?? false,
  },
  {
    title: '链接',
    icon: '🔗',
    action: () => openLinkDialog(),
    isActive: (e) => e?.isActive('link') ?? false,
  },
  {
    title: '图片',
    icon: '🖼',
    action: () => openImageUpload(),
  },
  {
    title: '表格',
    icon: '⊞',
    action: (_e, event) => openTablePicker(event),
  },
  {
    title: 'YouTube',
    icon: '▶',
    action: () => openYoutubeDialog(),
  },
  {
    title: '文字颜色',
    icon: '🎨',
    action: (e) => {
      const color = window.prompt('文字颜色 (例如 #e74c3c 或 red)');
      if (color) e?.chain().focus().setColor(color).run();
    },
  },
  {
    title: '高亮',
    icon: '🖍',
    action: (e) => e?.chain().focus().toggleHighlight().run(),
    isActive: (e) => e?.isActive('highlight') ?? false,
  },
  {
    title: '左对齐',
    icon: '⫷',
    action: (e) => e?.chain().focus().setTextAlign('left').run(),
    isActive: (e) => e?.isActive({ textAlign: 'left' }) ?? false,
  },
  {
    title: '居中',
    icon: '≣',
    action: (e) => e?.chain().focus().setTextAlign('center').run(),
    isActive: (e) => e?.isActive({ textAlign: 'center' }) ?? false,
  },
  {
    title: '右对齐',
    icon: '⫸',
    action: (e) => e?.chain().focus().setTextAlign('right').run(),
    isActive: (e) => e?.isActive({ textAlign: 'right' }) ?? false,
  },
  {
    title: '撤销',
    icon: '↩',
    action: (e) => e?.chain().focus().undo().run(),
  },
  {
    title: '重做',
    icon: '↪',
    action: (e) => e?.chain().focus().redo().run(),
  },
];
</script>

<style scoped>
.rich-editor {
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  background: var(--el-bg-color);
}

.editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-light);
}

.toolbar-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  color: var(--el-text-color-regular);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}

.toolbar-btn:hover {
  background: var(--el-fill-color);
  border-color: var(--el-border-color);
}

.toolbar-btn.active {
  background: var(--el-color-primary-light-9);
  border-color: var(--el-color-primary-light-5);
  color: var(--el-color-primary);
}

.toolbar-divider {
  display: inline-block;
  width: 1px;
  height: 22px;
  margin: 4px 4px 0;
  background: var(--el-border-color);
  align-self: center;
}

.editor-content {
  padding: 12px 16px;
  min-height: 240px;
  max-height: 480px;
  overflow-y: auto;
  font-size: 14px;
  line-height: 1.75;
  color: var(--el-text-color-regular);
}

.editor-content :deep(.ProseMirror) {
  outline: none;
  min-height: 220px;
}

.editor-content :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  color: var(--el-text-color-placeholder);
  pointer-events: none;
  float: left;
  height: 0;
}

.editor-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.editor-content :deep(th),
.editor-content :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 8px 12px;
  text-align: left;
  min-width: 60px;
}

.editor-content :deep(th) {
  background: var(--el-fill-color-light);
  font-weight: 600;
}

.editor-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.editor-content :deep(blockquote) {
  border-left: 3px solid var(--el-border-color-dark);
  padding-left: 16px;
  color: var(--el-text-color-secondary);
}

.editor-content :deep(pre) {
  background: #1e1e1e;
  color: #d4d4d4;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
}

.editor-content :deep(pre code) {
  background: none;
  color: inherit;
  font-size: 13px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
}

.editor-source {
  width: 100%;
  min-height: 320px;
  padding: 12px 16px;
  border: none;
  outline: none;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
  background: var(--el-bg-color);
  resize: vertical;
}

.table-picker-popup {
  position: fixed;
  z-index: 3000;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  padding: 8px;
  box-shadow: var(--el-box-shadow-light);
}

.table-picker-grid {
  display: grid;
  grid-template-columns: repeat(6, 24px);
  grid-template-rows: repeat(6, 24px);
  gap: 2px;
}

.grid-cell {
  width: 24px;
  height: 24px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 2px;
  cursor: pointer;
}

.grid-cell.selected {
  background: var(--el-color-primary-light-7);
  border-color: var(--el-color-primary-light-3);
}

.table-picker-label {
  text-align: center;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}
</style>
