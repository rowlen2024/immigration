export const usePermissions = () => {
  const permissions = useState<string[]>('admin-permissions', () => []);
  const loaded = useState<boolean>('admin-permissions-loaded', () => false);

  const loadPermissions = async () => {
    try {
      const api = useApi();
      const data = await api<{ permissions: string[] }>('/admin/me/permissions');
      permissions.value = data?.permissions ?? [];
      loaded.value = true;
    } catch {
      permissions.value = [];
      loaded.value = false;
    }
  };

  const hasPermission = (permission: string) => permissions.value.includes(permission);
  const hasAnyPermission = (required?: string[]) => {
    if (!required || required.length === 0) return true;
    return required.some((permission) => permissions.value.includes(permission));
  };

  const clearPermissions = () => {
    permissions.value = [];
    loaded.value = false;
  };

  return {
    permissions,
    loaded,
    loadPermissions,
    hasPermission,
    hasAnyPermission,
    clearPermissions,
  };
};
