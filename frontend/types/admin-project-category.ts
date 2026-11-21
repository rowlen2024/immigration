export type ProjectCategoryStatus = 'draft' | 'published';

export interface AdminProjectCategoryOption {
  id: number;
  name: string;
  slug: string;
  status: ProjectCategoryStatus;
}
