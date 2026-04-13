export const YJS_CONTENT_FIELD = 'content';
export const YJS_COMMENT_META_FIELD = 'comment_meta';

export const YJS_CONTENT_CARRIER = Object.freeze({
  TEXT: 'text',
  XML_FRAGMENT: 'xml-fragment'
});

const createTextBridge = (doc, field) => {
  const target = doc.getText(field);
  const toString = () => target.toString();

  return {
    target,
    isEmpty: () => target.length === 0,
    toString,
    replace: (value = '') => {
      target.delete(0, target.length);
      if (value) {
        target.insert(0, value);
      }
    },
    insertIfEmpty: (value = '') => {
      if (!value || target.length > 0) {
        return;
      }
      target.insert(0, value);
    },
    observeRemote: (onRemoteChange) => {
      const observer = (event) => {
        if (event?.transaction?.local) {
          return;
        }
        onRemoteChange?.(toString());
      };
      target.observe(observer);
      return () => target.unobserve(observer);
    }
  };
};

const createXmlFragmentBridge = (doc, field) => {
  const target = doc.getXmlFragment(field);
  return {
    target,
    isEmpty: () => {
      if (typeof target.length === 'number' && target.length > 0) {
        return false;
      }
      try {
        const json = target.toJSON?.();
        return !Array.isArray(json) || json.length === 0;
      } catch (_) {
        return true;
      }
    }
  };
};

const bridgeFactories = {
  [YJS_CONTENT_CARRIER.TEXT]: createTextBridge,
  [YJS_CONTENT_CARRIER.XML_FRAGMENT]: createXmlFragmentBridge
};

export const createYjsContentBridge = ({
  doc,
  carrier = YJS_CONTENT_CARRIER.TEXT,
  field = YJS_CONTENT_FIELD
} = {}) => {
  const factory = bridgeFactories[carrier] || bridgeFactories[YJS_CONTENT_CARRIER.TEXT];
  return factory(doc, field);
};
