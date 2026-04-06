import { computed } from 'vue';

export function useRichTextParagraphActions({
  editorRef,
  componentToolbarItemsRef,
  runComponentCommand,
  toggleBulletList,
  toggleOrderedList,
  toggleBlockquote,
  toggleCodeBlock,
  undo,
  redo
}) {
  const focusForParagraphAction = (ctx, position = 'start') => {
    if (typeof ctx?.focusBlock === 'function') {
      ctx.focusBlock(position);
      return;
    }
    editorRef.value?.commands?.focus?.();
  };

  const insertComponentByHandle = (ctx, item) => {
    if (!item) {
      return;
    }
    focusForParagraphAction(ctx);
    runComponentCommand(item);
  };

  const emptyParagraphActions = computed(() => {
    const actions = [
      {
        key: 'set-paragraph',
        label: '普通文本',
        iconKey: 'paragraph',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          editorRef.value?.chain().focus().setParagraph().run();
        }
      },
      {
        key: 'heading-1',
        label: '一级标题',
        iconKey: 'heading-1',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          editorRef.value?.chain().focus().setHeading({ level: 1 }).run();
        }
      },
      {
        key: 'heading-2',
        label: '二级标题',
        iconKey: 'heading-2',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          editorRef.value?.chain().focus().setHeading({ level: 2 }).run();
        }
      },
      {
        key: 'bullet-list',
        label: '无序列表',
        iconKey: 'bullet-list',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          toggleBulletList();
        }
      },
      {
        key: 'ordered-list',
        label: '有序列表',
        iconKey: 'ordered-list',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          toggleOrderedList();
        }
      },
      {
        key: 'blockquote',
        label: '引用',
        iconKey: 'blockquote',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          toggleBlockquote();
        }
      }
    ];

    const insertActions = [
      {
        key: 'code-block',
        label: '代码块',
        iconKey: 'code-block',
        handler: (ctx) => {
          focusForParagraphAction(ctx);
          toggleCodeBlock();
        }
      }
    ];

    componentToolbarItemsRef.value.forEach((item) => {
      insertActions.push({
        key: `component-${item.key}`,
        label: item.label,
        iconKey: item.key,
        handler: (ctx) => {
          insertComponentByHandle(ctx, item);
        }
      });
    });

    if (insertActions.length > 0) {
      actions.push({
        key: 'insert',
        label: '插入',
        iconKey: 'insert',
        children: insertActions
      });
    }

    actions.push({
      key: 'undo',
      label: '撤销',
      handler: () => {
        undo();
      }
    });

    actions.push({
      key: 'redo',
      label: '重做',
      handler: () => {
        redo();
      }
    });

    return actions;
  });

  return {
    emptyParagraphActions
  };
}
