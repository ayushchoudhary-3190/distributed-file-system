export interface FileItem {
  id: string;
  name: string;
  type: 'folder' | 'file';
  mimeType?: string;
  size?: number;
  modifiedAt: Date;
  owner: string;
  shared?: boolean;
  starred?: boolean;
  thumbnail?: string;
}

export type ViewMode = 'grid' | 'list';
export type SortField = 'name' | 'modified' | 'size';
export type SortOrder = 'asc' | 'desc';
