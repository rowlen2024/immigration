interface NavItem {
  id: number;
  label: string;
  link: string | null;
  children: NavItem[];
  status: boolean;
}

interface BreadcrumbItem {
  label: string;
  link?: string;
}

type NavPosition = 'header' | 'footer'

const FALLBACK_HEADER: NavItem[] = [
  {
    id: 1, label: '美国EB-5', link: '/projects/eb5', status: true,
    children: [
      { id: 6, label: 'EB-5项目概述', link: '/projects/eb5', children: [], status: true },
      { id: 7, label: '申请条件', link: '/projects/eb5#requirements', children: [], status: true },
      { id: 8, label: '申请流程', link: '/projects/eb5#timeline', children: [], status: true },
      { id: 9, label: '费用明细', link: '/projects/eb5#cost', children: [], status: true },
    ],
  },
  {
    id: 2, label: '香港投资', link: '/projects/cies', status: true,
    children: [
      { id: 10, label: '香港投资移民概述', link: '/projects/cies', children: [], status: true },
      { id: 11, label: '资格要求', link: '/projects/cies#requirements', children: [], status: true },
      { id: 12, label: '申请流程', link: '/projects/cies#timeline', children: [], status: true },
      { id: 13, label: '费用明细', link: '/projects/cies#cost', children: [], status: true },
    ],
  },
  {
    id: 3, label: '巴拿马购房', link: '/projects/panama', status: true,
    children: [
      { id: 14, label: '巴拿马购房移民概述', link: '/projects/panama', children: [], status: true },
      { id: 15, label: '购房要求', link: '/projects/panama#requirements', children: [], status: true },
      { id: 16, label: '申请流程', link: '/projects/panama#timeline', children: [], status: true },
      { id: 17, label: '费用明细', link: '/projects/panama#cost', children: [], status: true },
    ],
  },
  {
    id: 4, label: '项目对比', link: '/compare', children: [], status: true,
  },
  {
    id: 5, label: '关于我们', link: '/about', status: true,
    children: [
      { id: 18, label: '公司简介', link: '/about', children: [], status: true },
      { id: 19, label: '成功案例', link: '/cases', children: [], status: true },
      { id: 20, label: '常见问题', link: '/faq', children: [], status: true },
      { id: 21, label: '联系我们', link: '/contact', children: [], status: true },
    ],
  },
];

const FALLBACK_FOOTER: NavItem[] = [
  {
    id: 1, label: '移民项目', link: '', status: true,
    children: [
      { id: 11, label: '美国EB-5投资移民', link: '/projects/eb5', children: [], status: true },
      { id: 12, label: '香港投资移民', link: '/projects/cies', children: [], status: true },
      { id: 13, label: '巴拿马购房移民', link: '/projects/panama', children: [], status: true },
      { id: 14, label: '项目对比', link: '/compare', children: [], status: true },
    ],
  },
  {
    id: 2, label: '关于我们', link: '', status: true,
    children: [
      { id: 21, label: '公司简介', link: '/about', children: [], status: true },
      { id: 22, label: '成功案例', link: '/cases', children: [], status: true },
      { id: 23, label: '常见问题', link: '/faq', children: [], status: true },
      { id: 24, label: '联系我们', link: '/contact', children: [], status: true },
    ],
  },
];

const findNode = (items: NavItem[], path: string): NavItem | null => {
  for (const item of items) {
    if (item.link === path) return item;
    const found = findNode(item.children, path);
    if (found) return found;
  }
  return null;
};

const collectAncestors = (
  items: NavItem[],
  target: NavItem,
  ancestors: NavItem[] = [],
): NavItem[] | null => {
  for (const item of items) {
    if (item === target) return ancestors;
    const result = collectAncestors(item.children, target, [...ancestors, item]);
    if (result) return result;
  }
  return null;
};

const buildBreadcrumb = (navItems: NavItem[], path: string, label?: string): BreadcrumbItem[] => {
  const cleanPath = path.replace(/#.*$/, '');

  const node = findNode(navItems, cleanPath);

  if (node) {
    const ancestors = collectAncestors(navItems, node) || [];
    const items: BreadcrumbItem[] = [];
    for (const a of ancestors) {
      items.push({ label: a.label, link: a.link || undefined });
    }
    items.push({ label: node.label });

    if (label && label !== node.label) {
      items.push({ label });
    }
    return items;
  }

  // Try progressively shorter paths for unmatched routes (e.g. /compare/eb5-vs-cies)
  const segments = cleanPath.split('/').filter(Boolean);
  while (segments.length > 0) {
    segments.pop();
    const parentPath = '/' + segments.join('/');
    const parentResult = buildBreadcrumb(navItems, parentPath, undefined);
    if (parentResult.length > 0) {
      const lastSegment = path.split('/').filter(Boolean).pop() || '';
      parentResult.push({ label: label || lastSegment });
      return parentResult;
    }
  }

  const lastSeg = path.split('/').filter(Boolean).pop() || '';
  if (!lastSeg) return [];
  return [{ label: label || lastSeg }];
};

export const useNavigation = (position: NavPosition = 'header') => {
  const { data } = useFetch('/api/v1/navigation', {
    key: `navigation-${position}`,
    query: { position },
    transform: (response: any) => response?.data ?? response,
  })

  const fallback = position === 'footer' ? FALLBACK_FOOTER : FALLBACK_HEADER

  const navItems = computed<NavItem[]>(() => {
    const items = data.value as NavItem[] | null
    return (items && items.length > 0) ? items : fallback
  })

  const getBreadcrumb = (path: string, label?: string): BreadcrumbItem[] => {
    return buildBreadcrumb(navItems.value, path, label);
  };

  return { navItems, getBreadcrumb };
};
