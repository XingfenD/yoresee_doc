# 跨类型切换文档内容重复 / Cross-type document switch duplicates content

| 字段 | 值 |
| --- | --- |
| 状态 / Status | 已修复 / Fixed |
| 修复版本 / Fixed in | see `docs/CHANGELOG` |
| 影响范围 / Impact | 富文本编辑器 (`yoresee_rich_text`) ↔ 表格 (`yoresee_table`) ↔ 幻灯片 (`yoresee_slide`) 之间来回切换；不涉及 markdown（默认类型） |
| 数据后果 / Data consequence | YJS XmlFragment 内容被 CRDT 合并污染，snapshot worker 重新写回后端 `Document.Content`，**重复内容会被持久化** |
| 复现路径 / Reproduction | 富文本 doc A → 切到表格 doc B → 切回富文本 doc A → 看到重复段落 |

---

## 现象 / Symptom

用户在 YoreseeRichText 文档中正常编辑；切到表格文档做点操作，再切回原来的富文本文档，内容出现重复。

```
远思文档  <drawioblock diagram="              ">
测试
测试
远思
文档
测试测试
远思文档<drawioblock diagram="              ">
```

重复内容会被写入后端，刷新页面后仍在。

## 排查过程 / Investigation

### Phase 1 — 代码定位

锚定到三个组件：

- `frontend/src/views/document/DocumentEditor.vue:286` — 共享 `editorContent = ref('')`
- `frontend/src/composables/document/editor/useDocumentEditorLifecycle.js:56-73` — docId 切换时清空 `editorContent`
- `frontend/src/composables/document/editor/shared/useTypedDocumentPersistence.js:91-108` — table / slide 持久化通过 REST 拉内容

### Phase 2 — 时序假设

注意到 `useDocumentEditorLifecycle.js` 的 `docId` watcher 是 async，且含多个 `await`。在 `await` 让出主线程时，Vue 会 flush 微任务，触发 `useTypedDocumentPersistence` 的 `[docId, currentDocType]` watcher。在 `updateCurrentDocTitle()` 把 `currentDocType` 切到新类型之前，`isCurrentType` 短暂为 `true`（沿用旧类型），所以 `loadContent()` 会执行。

具体路径：从表格切到富文当时：
1. `docId` 变 → `editorContent = ''`（被生命周期清空）
2. await 让出 → `useTypedDocumentPersistence`（type='2' 或 '3'）watcher 触发，此时 `currentDocType` 还停在旧值，**`isCurrentType` 为 true**
3. `loadContent()` 通过 `getDocumentContent(newDocId)` 拉取新 doc 内容，写入共享 `editorContent`
4. 后端返回的是富文本 doc 的 `Content` 字段 = `Y.XmlFragment.toString()`（XML 序列化），被塞进共享 ref
5. `updateCurrentDocTitle()` 执行，`currentDocType` 切到富文本，富文本编辑器挂载，读到的 `editorContent` 已经是 XML
6. `markdownToHtml(xml)` → TipTap setContent → Collaboration extension 把这次 setContent 当作本地编辑写回 YJS XmlFragment（CRDT insert）
7. provider 同步服务端 YJS 状态，CRDT 合并「XML insert」+「服务端内容」→ 重复段落

### Phase 3 — 浏览器端 console 验证

为不依赖理论推导，加临时 `console.log` 在关键位置打印 `docId` / `currentDocType` / `isCurrentType` / `contentRef` / `responseDocType` / `modelValueRef.preview` / `YJS XmlFragment.toString()`。日志原文见附录 A。

关键证据：

```
[DEBUG][typed-persist] watcher fired { persistenceType: '3', docId: '94422bc...', currentDocType: '3', isCurrentType: true }
[DEBUG][typed-persist] loadContent response { responseDocType: '4', contentPreview: '<paragraph>ce shi</paragraph>...' }
...
[DEBUG][rich-text] mount { length: 0, preview: '' }
... loadContent 返回后 ...
[rich-text] yjs xmlfragment after sync { length: 4, toString: '<paragraph>ce shi</paragraph><paragraph>ce shice s…graph>ce shi</paragraph><paragraph>测试</paragraph>' }
```

确认了三件事：

1. **跨类型误触发**：slide persistence（`persistenceType: '3'`）拉到了 `responseDocType: '4'`（富文本）的文档
2. **写入被富文本编辑器消费**：富文本编辑器 `mount` 时 `preview: ''`（REST 还没回来），稍后 `setContent(xml)` 把 XML 内容塞进 TipTap，再被 Collaboration extension 写回 YJS
3. **CRDT 合并产生重复**：原本 `ce shi / ce shi / 测试` 三段，合并后变 4 段并出现串接字段 `ce shice s…`，snapshot worker 立刻把 `XmlFragment.toString()` 写回后端 → 持久化

### Phase 4 — 为什么 markdown ↔ rich-text 不出问题

markdown 与 rich-text 编辑器都不走 `useTypedDocumentPersistence`（只有 table / slide 走）。它们的 source of truth 是 YJS XmlFragment / Y.Text；`editorContent` 只作为 v-model 的双向桥。所以表格/幻灯片 persistence 跨类型误触发不会污染 markdown/rich-text 的实际存储。

## 修复 / Fix

### 选定方案：按类型拆分 editorContent（方案 C）

三个候选方案：

| 方案 | 描述 | 取舍 |
| --- | --- | --- |
| A：loadContent 后置校验 `doc.type` | 拿到响应后比对 `response.document.type` 与 persistence 自己的 `type`，不匹配直接 return | 最小改动，仍发了一次无效 HTTP 请求 |
| B：watcher 前置去抖 + 类型校验 | `await nextTick()` 让 currentDocType 落定后再判断 | 依赖 lifecycle 的 await 数量，脆弱耦合 |
| C：editorContent 按类型拆分 | 每个类型持有自己的 content ref，persistence 写自己的，编辑器读自己的，**互不共享** | 改动面较大但架构上最干净，根除共享状态被污染的可能 |

最终选 **C**：从根上消除「共享 content ref 被跨类型持久化污染」的可能。代价是改 3 个文件、46 行新增 / 15 行删除。

### 改动清单

`frontend/src/views/document/DocumentEditor.vue`
- 新增 `markdownContent` / `tableContent` / `slideContent` / `richTextContent` 四个 ref
- 每个编辑器 `v-model` 绑到对应的 ref
- `useTableDocumentPersistence` 传 `tableContent`，`useSlideDocumentPersistence` 传 `slideContent`
- `useDocumentEditorActions` / `useDocumentEditorLifecycle` 都改为接收四个 ref

`frontend/src/composables/document/editor/useDocumentEditorLifecycle.js`
- 新增 `clearEditorContents()`，在 docId watcher 里一次性清空四个 ref（之前只清一个共享 ref）

`frontend/src/composables/document/editor/useDocumentEditorActions.js`
- 新增 `resolveActiveContent()`：根据 `currentDocType` 选对应的 ref，模板创建（`submitCreateTemplate` / `openCreateTemplateDialog`）从这里读

### 修复效果（再次用 console probe 验证）

修复后跨类型切换的日志：

```
[DEBUG][typed-persist] loadContent called { context: 'loadSlideDocument', persistenceType: '3', contentRef: 'slideContent', docId: '94422bc...', currentDocType: '3', isCurrentType: true }
[DEBUG][typed-persist] loadContent response { context: 'loadSlideDocument', contentRef: 'slideContent', responseDocType: '4', contentPreview: '<paragraph>测试</paragraph>' }
[DEBUG][typed-persist] loadContent DONE { contentRef: 'slideContent', docId: '94422bc...', len: 25 }
```

注意 **`contentRef: 'slideContent'`** —— 即使 persistence 误触发拉到了富文本 doc 的 XML 内容，它写入的是 `slideContent`，而当前显示的是富文本编辑器（读 `richTextContent`）。两个 ref 完全隔离。

用户视角：「没有重复了」。

## 为什么之前没人发现

- markdown ↔ rich-text 切换是最高频路径（不触发持久化），单元测试只覆盖了这路径
- 跨类型 bug 需要「来回切」才能暴露，单次切换看起来正常
- 富文本内容被 snapshot worker 写回后端时已经包含 XML 序列化，正常用户看到的内容 vs 被污染的内容差距不大（TipTap 会尽量解析 `<paragraph>` 等标签，但 schema 不识别就丢弃），复现门槛高

## 复盘 / Lessons

1. **共享 ref 是跨类型持久化的天然雷区** —— 任何「单一容器承载多种类型」的架构都该意识到：类型切换中间态存在时序窗口，持有者必须基于最终态决策而非瞬时态
2. **后端 Content 字段对富文本存的是 YJS XmlFragment 的 XML 序列化** —— 不能直接当 markdown 用；这次 bug 的放大器正是 `markdownToHtml(xml-string)` 把 XML 当 markdown 解析失败后让 Collaboration extension 写回 YJS
3. **CRDT 不是万能的** —— 它保证 eventual consistency，但当「本地初始状态」和「远端同步内容」都被插入到 XmlFragment 时，两条 INSERT 路径都会被保留，导致重复。修这种 bug 要在源头上避免无效 INSERT，而不是依赖 CRDT 去重

## 附录 A：验证用的临时 console.log 位置（已全部移除）

| 文件 | 用途 |
| --- | --- |
| `frontend/src/composables/document/editor/useDocumentEditorLifecycle.js:58,69` | 标记 docId watcher 触发和 `updateCurrentDocTitle` 完成 |
| `frontend/src/composables/document/editor/shared/useTypedDocumentPersistence.js:92,109,123` | 标记 persistence watcher / loadContent / response / DONE |
| `frontend/src/composables/document/editor/rich-text-editor/useRichTextEditorRuntime.js:177` | 富文本 mount 时的 modelValueRef 长度与前 200 字 |
| `frontend/src/composables/document/editor/rich-text-editor/useRichTextYjsCollaboration.js:122,127` | YJS provider sync 完成时 XmlFragment 长度和 toString |

修复并验证后已全部移除；`window.__markdownContent` 等调试用 ref 暴露也已清理。
