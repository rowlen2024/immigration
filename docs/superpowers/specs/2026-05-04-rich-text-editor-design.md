# Rich Text Editor for Page Content Form

## Summary

Replace the plain `<el-input type="textarea">` in the admin page editor with a Tiptap rich text editor, supporting full formatting (tables, images, video, syntax highlight, color) and a source-code toggle for raw HTML editing.

## Scope

- Admin: new/edit page drawer (pages.vue)
- New shared RichEditor component
- Backend sanitization policy update (page_svc.go)
- No changes to public display ([...slug].vue already renders v-html)

## Frontend

### Dependencies (14 new packages)

@tiptap/vue-3, @tiptap/starter-kit, @tiptap/extension-table, @tiptap/extension-table-row, @tiptap/extension-table-cell, @tiptap/extension-table-header, @tiptap/extension-image, @tiptap/extension-text-style, @tiptap/extension-color, @tiptap/extension-highlight, @tiptap/extension-text-align, @tiptap/extension-youtube, @tiptap/extension-underline, @tiptap/extension-code-block-lowlight

### New: `frontend/components/RichEditor.vue`

Encapsulated Tiptap editor with:

- **Toolbar**: bold, italic, underline, heading 1-3, blockquote, ordered/unordered list, link, image, table, youtube embed, code block, text color, highlight, text-align (left/center/right), undo, redo, source toggle
- **v-model** binding: emits HTML string, accepts HTML string
- **Image upload**: on file select, POST to `/api/v1/admin/media/upload` with Bearer token, insert returned URL into editor
- **Source toggle**: button switches between WYSIWYG and raw HTML `<textarea>`; round-trips through `editor.getHTML()` / `editor.commands.setContent()`

Props: `modelValue: string` (HTML content).
Emits: `update:modelValue`.

### Modify: `frontend/pages/admin/pages.vue`

- Replace `<el-input v-model="form.content" type="textarea" :rows="8" />` with `<RichEditor v-model="form.content" />`
- Drawer width: `560px` → `900px`
- No other changes (validation rules, save logic unchanged)

## Backend

### Modify: `backend/internal/service/page_svc.go`

Replace `bluemonday.UGCPolicy()` with a custom policy that allows the rich HTML tags Tiptap produces:

- **Elements**: h1-h6, p, br, hr, ul, ol, li, blockquote, pre, code, strong, em, u, s, del, a, img, table, thead, tbody, tr, th, td, caption, div, span, iframe, video, source
- **Attributes**: src/alt/title/width/height on img; href/title/target/rel on a; src/frameborder/allowfullscreen on iframe; src/controls/width/height on video/source; style/class on span/div/td/th for text styling and alignment
- **Styles**: color, background-color, text-align on span/td/th
- **URL schemes**: http, https, mailto only
- **Links**: require rel=nofollow

## Verification

- Create a page with all formatting types (table, image, video embed, colored text, code block, alignment) → save → edit → verify content preserved
- Paste HTML source → toggle to WYSIWYG → toggle back → verify HTML intact
- Upload an image via the editor toolbar → verify it appears in media library
- Open the public page → verify rich content renders correctly with existing styles
- Submit malicious HTML (onerror, javascript:) → verify sanitization strips it
