import { defineStore } from 'pinia';

interface Project {
  id: string;
  title: string;
  slug: string;
  summary: string;
  description: string;
  coverImage: string;
  category: string;
  minInvestment: string;
  processingTime: string;
  requirements: string[];
  benefits: string[];
  status: string;
  createdAt: string;
  updatedAt: string;
}

interface ProjectsState {
  projects: Project[];
  currentProject: Project | null;
  loading: boolean;
  error: string | null;
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectsState => ({
    projects: [],
    currentProject: null,
    loading: false,
    error: null,
  }),

  getters: {
    activeProjects: (state) =>
      state.projects.filter((p) => p.status === 'published'),

    projectBySlug: (state) => (slug: string) =>
      state.projects.find((p) => p.slug === slug) || null,
  },

  actions: {
    async fetchProjects() {
      this.loading = true;
      this.error = null;

      try {
        const api = useApi();
        const data = await api<Project[]>('/projects');
        this.projects = data;
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch projects';
      } finally {
        this.loading = false;
      }
    },

    async fetchProject(slug: string) {
      this.loading = true;
      this.error = null;

      try {
        const api = useApi();
        const data = await api<Project>(`/projects/${slug}`);
        this.currentProject = data;
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Failed to fetch project';
      } finally {
        this.loading = false;
      }
    },
  },
});
