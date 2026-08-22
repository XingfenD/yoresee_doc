export const inviteStatusType = (status) => {
  if (status === 'active') return 'success';
  if (status === 'expired') return 'info';
  return 'warning';
};

export const inviteStatusLabel = (status, t) => {
  if (status === 'active') return t('user.invite.active');
  if (status === 'expired') return t('user.invite.expired');
  return t('user.invite.disabled');
};

const pad = (num) => `${num}`.padStart(2, '0');

export const formatDateYYYYMMDD = (date) =>
  `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;

export const toRfc3339EndOfDay = (dateStr) => {
  if (!dateStr) return '';
  if (dateStr.includes('T')) return dateStr;
  const [year, month, day] = dateStr.split('-').map((value) => Number(value));
  if (!year || !month || !day) return '';
  const local = new Date(year, month - 1, day, 23, 59, 59);
  return local.toISOString();
};

export const resolveInviteStatus = (invite) => {
  if (invite.disabled) return 'disabled';
  if (invite.expires_at) {
    const expireTs = Date.parse(invite.expires_at);
    if (!Number.isNaN(expireTs) && expireTs < Date.now()) return 'expired';
  }
  if (
    typeof invite.max_used_cnt === 'number' &&
    typeof invite.used_cnt === 'number' &&
    invite.used_cnt >= invite.max_used_cnt
  ) {
    return 'expired';
  }
  return 'active';
};

export const toNumber = (value, fallback = null) => {
  if (value === null || value === undefined) return fallback;
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
};

export const mapInviteRow = (invite, { isSystemMode, t }) => {
  const row = {
    code: invite.code,
    created_at: invite.created_at || '-',
    expires_at: invite.expires_at || '-',
    status: resolveInviteStatus(invite),
    max: toNumber(invite.max_used_cnt, null),
    used: toNumber(invite.used_cnt, 0),
    disabled: invite.disabled,
    note: invite.note || ''
  };
  if (isSystemMode) {
    row.created_by = invite.created_by_name || invite.created_by_external_id || t('common.unknown');
  }
  return row;
};

export const mapRecordRow = (record) => ({
  row_key: record.row_key || `${record.code || ''}_${record.used_at || ''}_${record.status || ''}`,
  code: record.code,
  used_by: record.used_by || '-',
  used_at: record.used_at || '-',
  status: record.status || 'failed'
});
